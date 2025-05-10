package main

import (
	"io"
	"log/slog"
	"net"
)

func main() {
	slog.Info("starting")

	l, err := net.Listen("tcp", "0.0.0.0:25565")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	upstream, err := net.Dial("tcp", "127.0.0.1:15550")
	if err != nil {
		panic(err)
	}
	defer upstream.Close()

	done := make(chan bool)

	buf := make([]byte, 4096)
	go func() {
		io.CopyBuffer(upstream, c, buf)
	}()

	buf1 := make([]byte, 4096)
	go func() {
		io.CopyBuffer(c, upstream, buf1)
	}()

	<-done
}
