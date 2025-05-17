package handler

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
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
	// TODO send control frame
	c.conn.Close()
}

func (c *WebSocketConnection) Read(p []byte) (n int, err error) {
	var buf bytes.Buffer

	fin := false

	for !fin {
		b0, _ := c.rw.ReadByte()
		fin = b0&0b10000000 != 0  // fin
		opcode := b0 & 0b00001111 // opcode

		// TODO handle correclty
		if opcode == WebSocketOpcodeConnectionCloseFrame {
			return n, io.EOF
		}

		b1, _ := c.rw.ReadByte()
		mask := b1&0b10000000 != 0
		plen := b1 & 0b01111111

		var uintplen uint64
		if plen <= 125 {
			uintplen = uint64(plen)
		} else if plen == 126 {
			x := make([]byte, 2)
			io.ReadFull(c.rw.Reader, x)
			uintplen = uint64(binary.BigEndian.Uint16(x))
		} else if plen == 127 {
			x := make([]byte, 8)
			io.ReadFull(c.rw.Reader, x)
			uintplen = binary.BigEndian.Uint64(x)
		}

		if !mask {
			panic("websocket read: mask must be set")
		}

		maskingKey := make([]byte, 4)
		io.ReadFull(c.rw.Reader, maskingKey)

		payload := make([]byte, uintplen)
		io.ReadFull(c.rw.Reader, payload)

		msgPart := make([]byte, uintplen)
		for i := 0; i < int(uintplen); i++ {
			msgPart[i] = payload[i] ^ maskingKey[i%4]
		}

		buf.Write(p)
	}

	return buf.Read(p)
}

func (c *WebSocketConnection) Write(p []byte) (n int, err error) {
	// TODO use 16KB payload??
	var buf bytes.Buffer

	buf.WriteByte(0b10000001)

	plen := len(p)

	switch {
	case plen <= 125:
		buf.WriteByte(byte(plen))
	case plen <= 65535:
		buf.WriteByte(126)
		binary.Write(&buf, binary.BigEndian, uint16(plen))
	default:
		buf.WriteByte(127)
		binary.Write(&buf, binary.BigEndian, uint64(plen))
	}

	buf.Write(p)

	_, err = c.rw.Write(buf.Bytes())
	if err != nil {
		return 0, err
	}

	err = c.rw.Flush()
	if err != nil {
		return 0, err
	}

	return plen, nil
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
