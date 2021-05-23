package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type (
	RESPCommand      string
	CommandDelimiter string
)

const (
	RESPCommand_UNKNOWN RESPCommand = "UNKNOWN"
	RESPCommand_PING                = "PING"
	RESPCommand_ECHO                = "ECHO"

	CommandDelimiterSimpleStrings = "+"
	CommandDelimiterErrors        = "-"
	CommandDelimiterIntegers      = ":"
	CommandDelimiterBulkStrings   = "$"
	CommandDelimiterArrays        = "*"

	COMMAND_DELIMITER_LENGTH = 1

	CRLF = "\r\n"
)

var (
	ErrInvalidMessage = errors.New("invalid message")
)

type RESPDecoder struct {
	raw string
}
type RESPResult struct {
	Command RESPCommand
	Raw     []string
}

func NewRESPDecoder(reader io.Reader) (*RESPDecoder, error) {
	buf := make([]byte, 1024) // TODO support large size
	if _, err := reader.Read(buf); err != nil {
		return nil, err
	}

	str := string(buf)
	if !strings.Contains(str, "\r\n") {
		return nil, fmt.Errorf("message should have carriage returns")
	}

	return &RESPDecoder{
		raw: str,
	}, nil
}

func (d *RESPDecoder) Decode() (*RESPResult, error) {
	del, err := d.get(0, 1)
	if err != nil {
		return nil, err
	}

	var raw []string
	if del == CommandDelimiterArrays {
		if err := d.seek(1); err != nil {
			return nil, err
		}

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
			arr[i] = str
		}

		raw = arr
	} else {
		str, err := d.decode()
		if err != nil {
			return nil, err
		}
		raw = []string{str}
	}

	var cmd RESPCommand
	switch strings.ToUpper(raw[0]) {
	case RESPCommand_PING:
		cmd = RESPCommand_PING
	case RESPCommand_ECHO:
		cmd = RESPCommand_ECHO
	default:
		cmd = RESPCommand_UNKNOWN
	}

	return &RESPResult{
		Command: cmd,
		Raw:     raw,
	}, nil
}

func (d *RESPDecoder) get(start, end int) (string, error) {
	if start > len(d.raw) || end > len(d.raw) || start > end {
		return "", fmt.Errorf("failed to get bytes")
	}
	return d.raw[start:end], nil
}

func (d *RESPDecoder) seek(size int) error {
	if len(d.raw) < size {
		return fmt.Errorf("failed to seek message")
	}
	d.raw = d.raw[size:]
	return nil
}

func (d *RESPDecoder) seekToCRLF() (string, error) {
	arr := strings.Split(d.raw, CRLF)
	if len(arr) == 0 {
		return "", fmt.Errorf("failed to seek message")
	}
	el := arr[0]
	d.seek(len(el) + len(CRLF))
	return el, nil
}

func (d *RESPDecoder) decode() (string, error) {
	del, err := d.get(0, 1)
	if err != nil {
		return "", err
	}
	if err := d.seek(1); err != nil {
		return "", err
	}

	switch del {
	case CommandDelimiterSimpleStrings:
		str, err := d.decodeString()
		if err != nil {
			return "", fmt.Errorf("%w %v", ErrInvalidMessage, err)
		}
		return str, nil
	case CommandDelimiterIntegers:
		str, err := d.decodeInteger()
		if err != nil {
			return "", fmt.Errorf("%w %v", ErrInvalidMessage, err)
		}
		return str, nil
	case CommandDelimiterBulkStrings:
		str, err := d.decodeBulkStrings()
		if err != nil {
			return "", fmt.Errorf("%w %v", ErrInvalidMessage, err)
		}
		return str, nil
	default:
		return "", fmt.Errorf("%w invalid command delimiter. del=%s", ErrInvalidMessage, string(del))
	}
}

func (d *RESPDecoder) decodeString() (string, error) {
	str, err := d.seekToCRLF()
	if err != nil {
		return "", err
	}
	return str, nil
}

func (d *RESPDecoder) decodeInteger() (string, error) {
	str, err := d.seekToCRLF()
	if err != nil {
		return "", err
	}
	if _, err := strconv.Atoi(str); err != nil {
		return "", err
	}
	return str, nil
}

func (d *RESPDecoder) decodeBulkStrings() (string, error) {
	sizeStr, err := d.seekToCRLF()
	if err != nil {
		return "", err
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return "", err
	}
	str, err := d.seekToCRLF()
	if err != nil {
		return "", err
	}
	if len(str) != size {
		return "", fmt.Errorf("size should be equal to message length")
	}
	return str, nil
}
