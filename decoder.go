package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type (
	RESPCommand      int
	CommandDelimiter string
)

const (
	RESPCommand_UNKNOWN RESPCommand = iota
	RESPCommand_PING
	RESPCommand_ECHO

	CommandDelimiterSimpleStrings = "+"
	CommandDelimiterErrors        = "-"
	CommandDelimiterIntegers      = ":"
	CommandDelimiterBulkStrings   = "$"
	CommandDelimiterArrays        = "*"

	COMMAND_DELIMITER_LENGTH = 1
	CRLF                     = "\r\n"
)

var (
	ErrInvalidMessage = errors.New("invalid message")
)

type RESPDecoder struct {
	raw []byte
}

func NewRESPDecoder(reader io.Reader) (*RESPDecoder, error) {
	buf := make([]byte, 1024) // TODO support large size
	if _, err := reader.Read(buf); err != nil {
		return nil, err
	}

	return &RESPDecoder{
		raw: buf,
	}, nil
}

func (d *RESPDecoder) seek(size int) error {
	if len(d.raw) < size {
		return fmt.Errorf("failed to seek message")
	}
	d.raw = d.raw[size:]
	return nil
}

func (d *RESPDecoder) seekToCRLF() ([]byte, error) {
	arr := bytes.Split(d.raw, []byte(CRLF))
	if len(arr) == 0 {
		return nil, fmt.Errorf("failed to seek message")
	}
	el := arr[0]
	d.seek(len(string(el)) + len(CRLF))
	return el, nil
}

func (d *RESPDecoder) Decode() ([]string, error) {
	del := d.raw[0]
	if err := d.seek(1); err != nil {
		return nil, err
	}

	if string(del) == CommandDelimiterArrays {
		sizeStr, err := d.seekToCRLF()
		if err != nil {
			return nil, err
		}
		size, err := strconv.Atoi(string(sizeStr))
		if err != nil {
			return nil, err
		}

		arr := make([]string, size)
		for i := 0; i < size; i++ {
			str, err := d.decode()
			if err != nil {
				return nil, err
			}
			arr[i] = string(str)
		}
		return arr, nil
	}

	if _, err := d.seekToCRLF(); err != nil {
		return nil, err
	}
	str, err := d.decode()
	if err != nil {
		return nil, err
	}
	return []string{str}, nil
}

func (d *RESPDecoder) decode() (string, error) {
	del := d.raw[0]
	if err := d.seek(1); err != nil {
		return "", err
	}

	switch string(del) {
	case CommandDelimiterSimpleStrings:
		str, err := d.decodeString()
		if err != nil {
			return "", fmt.Errorf("%w %v", ErrInvalidMessage, err)
		}
		return str, nil
	case CommandDelimiterIntegers:
		// TODO
		return "", nil
	case CommandDelimiterBulkStrings:
		str, err := d.decodeBulkStrings()
		if err != nil {
			return "", fmt.Errorf("%w %v", ErrInvalidMessage, err)
		}
		return str, nil
	default:
		return "", fmt.Errorf("%w invalid command delimiter %s.", ErrInvalidMessage, string(del))
	}
}

func (d *RESPDecoder) decodeString() (string, error) {
	str, err := d.seekToCRLF()
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func (d *RESPDecoder) decodeInteger() (string, error) {
	str, err := d.seekToCRLF()
	if err != nil {
		return "", err
	}
	if _, err := strconv.Atoi(string(str)); err != nil {
		return "", err
	}
	return string(str), nil
}

func (d *RESPDecoder) decodeBulkStrings() (string, error) {
	sizeStr, err := d.seekToCRLF()
	if err != nil {
		return "", err
	}
	size, err := strconv.Atoi(string(sizeStr))
	if err != nil {
		return "", err
	}
	str, err := d.seekToCRLF()
	if err != nil {
		return "", err
	}
	if len(string(str)) != size {
		return "", err
	}
	return string(str), nil
}
