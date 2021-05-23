package main

import (
	"io"
	"log"
)

var (
	_ Handler = handler
)

func handler(conn io.ReadWriteCloser) error {
	defer conn.Close()

	req := make([]byte, 100)
	if _, err := conn.Read(req); err != nil {
		return err
	}

	log.Print(string(req))

	if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
		return err
	}
	return nil
}
