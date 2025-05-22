package handler

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
)

func ReadRequestJson(value io.ReadCloser, object interface{}) {
	d := json.NewDecoder(value)
	err := d.Decode(object)
	if err != nil {
		panic(err)
	}
}

func WriteResponseJson(w http.ResponseWriter, statusCode int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	c := json.NewEncoder(w)
	c.SetEscapeHTML(true)
	err := c.Encode(value)
	if err != nil {
		panic(err)
	}
}

const (
	WebSocketOpcodeContinuationFrame    = 0  // RFC 6455, 11.8
	WebSocketOpcodeTextFrame            = 1  // RFC 6455, 11.8
	WebSocketOpcodeBinaryFrame          = 2  // RFC 6455, 11.8
	WebSocketOpcodeConnectionCloseFrame = 8  // RFC 6455, 11.8
	WebSocketOpcodePingFrame            = 9  // RFC 6455, 11.8
	WebSocketOpcodePongFrame            = 10 // RFC 6455, 11.8
)

type WebSocketConnection struct {
	conn net.Conn
	rw   *bufio.ReadWriter
}

func NewWebSocketConnection(conn net.Conn, rw *bufio.ReadWriter) *WebSocketConnection {
	return &WebSocketConnection{conn: conn, rw: rw}
}

func (c *WebSocketConnection) Close() {
	if c.conn == nil {
		return
	}

	// TODO send control frame
	var msg []byte

	msg = append(msg, byte(0b10001000))
	msg = append(msg, byte(0b00000010))

	code := make([]byte, 2)
	binary.BigEndian.PutUint16(code, uint16(1000))
	msg = append(msg, code...)

	_, err := c.rw.Write(msg)
	if err != nil {
		panic(err)
	}

	err = c.rw.Flush()
	if err != nil {
		panic(err)
	}

	c.conn.Close()
}

func (c *WebSocketConnection) Read(p []byte) (n int, err error) {
	var fullPayload []byte

	b0, err := c.rw.ReadByte()
	if err != nil {
		return 0, err
	}
	fin := b0&0b10000000 != 0
	opcode := b0 & 0b00001111

	if opcode == WebSocketOpcodeConnectionCloseFrame {
		return 0, io.EOF
	}

	b1, err := c.rw.ReadByte()
	if err != nil {
		return 0, err
	}
	mask := b1&0b10000000 != 0
	if !mask {
		return 0, errors.New("websocket read: mask must be set")
	}

	payloadLength := b1 & 0b01111111
	var payloadLen uint64

	switch payloadLength {
	case 126:
		var lenBuf [2]byte
		if _, err := io.ReadFull(c.rw, lenBuf[:]); err != nil {
			return 0, err
		}
		payloadLen = uint64(binary.BigEndian.Uint16(lenBuf[:]))
	case 127:
		var lenBuf [8]byte
		if _, err := io.ReadFull(c.rw, lenBuf[:]); err != nil {
			return 0, err
		}
		payloadLen = binary.BigEndian.Uint64(lenBuf[:])
	default:
		payloadLen = uint64(payloadLength)
	}

	var maskingKey [4]byte
	if _, err := io.ReadFull(c.rw, maskingKey[:]); err != nil {
		return 0, err
	}

	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(c.rw, payload); err != nil {
		return 0, err
	}

	for i := range payload {
		payload[i] ^= maskingKey[i%4]
	}

	fullPayload = payload
	if !fin {
		// Note: you’re not handling fragmented frames here
		return 0, errors.New("fragmented frames not supported")
	}

	n = copy(p, fullPayload)
	return n, nil
}

func (c *WebSocketConnection) Write(p []byte) (n int, err error) {
	// TODO use 16KB payload??
	var msg []byte

	msg = append(msg, byte(0b10000001)) // fin, ..., opcode

	payloadLength := len(p)
	if payloadLength <= 125 {
		msg = append(msg, byte(payloadLength))
	} else if payloadLength <= 65535 {
		msg = append(msg, byte(126))
		x := make([]byte, 2)
		binary.BigEndian.PutUint16(x, uint16(payloadLength))
		msg = append(msg, x...)
	} else if payloadLength > 65535 {
		msg = append(msg, byte(127))
		x := make([]byte, 8)
		binary.BigEndian.PutUint64(x, uint64(payloadLength))
		msg = append(msg, x...)
	}

	// fmt.Print(payloadLength, string(p))
	msg = append(msg, p...) // payload

	n, err = c.rw.Write(msg)
	if err != nil {
		return n, err
	}

	err = c.rw.Flush()
	if err != nil {
		return n, err
	}

	return payloadLength, nil
}

func Upgrade(w http.ResponseWriter, r *http.Request) *WebSocketConnection {
	if r.Method != "GET" {
		panic("upgrade: method must be GET")
	}

	if r.Header.Get("Connection") == "" {
		panic("upgrade: missing Connection header")
	}

	if !strings.Contains(r.Header.Get("Connection"), "Upgrade") {
		panic("upgrade: missing Connection: Upgrade header")
	}

	secWebSocketKey := r.Header.Get("Sec-WebSocket-Key")
	if secWebSocketKey == "" {
		panic("upgrade: missing Sec-WebSocket-Key header")
	}

	if r.Header.Get("Sec-WebSocket-Version") != "13" {
		panic("upgrade: websocket version invalid")
	}

	h := sha1.New()
	_, err := h.Write([]byte(secWebSocketKey + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	if err != nil {
		panic(err)
	}
	secWebSocketAcceptHash := h.Sum(nil)
	secWebSocketAcceptHashEncoded := base64.StdEncoding.EncodeToString(secWebSocketAcceptHash)

	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	w.Header()["Sec-WebSocket-Accept"] = []string{secWebSocketAcceptHashEncoded}
	w.WriteHeader(http.StatusSwitchingProtocols)

	hj, ok := w.(http.Hijacker)
	if !ok {
		panic("response writer doesn't support hijacking")
	}
	netConn, rw, err := hj.Hijack()
	if err != nil {
		panic(err)
	}

	return NewWebSocketConnection(netConn, rw)
}
