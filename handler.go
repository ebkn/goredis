package main

import (
	"fmt"
	"io"
	"log"
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

	decoded, err := decoder.Decode()
	if err != nil {
		if _, err := conn.Write(formatError(err)); err != nil {
			return err
		}
		return err
	}

	encoder := NewRESPEncoder()

	switch decoded.Command {
	case RESPCommand_COMMAND:
		return nil
	case RESPCommand_PING:
		var str string
		switch len(decoded.Raw) {
		case 1:
			str = RESPResponse_PONG
		case 2:
			str = decoded.Raw[1]
		default:
			errMsg := fmt.Errorf("wrong number of arguments")
			if _, err := conn.Write(formatError(errMsg)); err != nil {
				return err
			}
			return nil
		}

		msg, err := encoder.EncodeBulkStrings(str)
		if err != nil {
			return err
		}
		if _, err := conn.Write(msg); err != nil {
			return err
		}
	case RESPCommand_ECHO:
		if len(decoded.Raw) != 2 {
			errMsg := fmt.Errorf("wrong number of arguments")
			if _, err := conn.Write(formatError(errMsg)); err != nil {
				log.Println(err)
				return nil
			}
			return nil
		}

		raw := decoded.Raw[1]
		msg, err := encoder.EncodeBulkStrings(raw)
		if err != nil {
			return err
		}
		if _, err := conn.Write(msg); err != nil {
			return err
		}
	default:
		errMsg := fmt.Errorf("unknown command %s", decoded.Raw[0])
		if _, err := conn.Write(formatError(errMsg)); err != nil {
			return err
		}
	}

	return nil
}

func formatError(err error) []byte {
	return []byte(fmt.Sprintf("%s%v%s", CommandDelimiterErrors, err, CRLF))
}
