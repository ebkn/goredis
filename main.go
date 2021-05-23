package main

import (
	"io"
	"log"
)

func main() {
	s := NewServer(6379)
	if err := s.Start(handler); err != nil {
		log.Fatal(err)
	}
}

func handler(rw io.ReadWriteCloser) error {
	defer rw.Close()

	req := make([]byte, 100)
	if _, err := rw.Read(req); err != nil {
		return err
	}

	log.Print(string(req))

	if _, err := rw.Write(req); err != nil {
		return err
	}
	return nil
}
