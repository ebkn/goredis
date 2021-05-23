package main

import (
	"log"
)

func main() {
	s := NewServer(6379)
	if err := s.Start(handler); err != nil {
		log.Fatal(err)
	}
}
