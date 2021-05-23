package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

type Handler func(r io.ReadWriteCloser) error

func (s *Server) Start(handler Handler) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	defer ln.Close()

	log.Printf("Starting server at :%d.\n", s.port)

	errChan := make(chan error, 1)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		select {
		case err := <-errChan:
			return err
		default:
			go func() {
				if err := handler(conn); err != nil {
					errChan <- err
				}
			}()
		}
	}
}
