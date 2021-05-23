package main

import (
	"bytes"
	"fmt"
)

type (
	RESPResponse string
)

const (
	RESPResponse_PONG = "PONG"
)

type RESPEncoder struct{}

func NewRESPEncoder() *RESPEncoder {
	return &RESPEncoder{}
}

func (e *RESPEncoder) EncodeString(str string) ([]byte, error) {
	msg := fmt.Sprintf(
		"%s%s%s",
		CommandDelimiterSimpleStrings,
		str,
		CRLF,
	)
	return []byte(msg), nil
}

func (e *RESPEncoder) EncodeBulkStrings(str string) ([]byte, error) {
	msg := fmt.Sprintf(
		"%s%d%s%s%s",
		CommandDelimiterBulkStrings,
		len(str),
		CRLF,
		str,
		CRLF,
	)
	return []byte(msg), nil
}

func (e *RESPEncoder) EncodeStringSlice(arr []string) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := buf.WriteString(fmt.Sprintf(
		"%s%d%s",
		CommandDelimiterArrays,
		len(arr),
		CRLF,
	)); err != nil {
		return nil, err
	}
	for _, el := range arr {
		msg, err := e.EncodeBulkStrings(el)
		if err != nil {
			return nil, err
		}
		if _, err := buf.Write(msg); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
