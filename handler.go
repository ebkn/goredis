package main

import (
	"fmt"
	"io"
)

var (
	_ Handler = handler
)

func handler(conn io.ReadWriteCloser) error {
	defer conn.Close()

	decoder, err := NewRESPDecoder(conn)
	if err != nil {
		return err
	}

	req, err := decoder.Decode()
	if err != nil {
		// TODO
		if _, err := conn.Write([]byte(err.Error())); err != nil {
			return err
		}
		return nil
	}

	fmt.Println(req)

	if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
		return err
	}
	return nil
}
