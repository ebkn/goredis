package main

import (
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

	decoded, err := decoder.Decode()
	if err != nil {
		// TODO
		if _, err := conn.Write([]byte(err.Error())); err != nil {
			return err
		}
		return nil
	}

	encoder := NewRESPEncoder()

	switch decoded.Command {
	case RESPCommand_PING:
		var str string
		switch len(decoded.Raw) {
		case 1:
			str = RESPResponse_PONG
		case 2:
			str = decoded.Raw[1]
		default:
			// TODO
			if _, err := conn.Write([]byte("wrong number of arguments")); err != nil {
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
			// TODO
			if _, err := conn.Write([]byte("wrong number of arguments")); err != nil {
				return err
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
		// TODO initial command
		msg, err := encoder.EncodeString(RESPResponse_PONG)
		if err != nil {
			return err
		}
		if _, err := conn.Write(msg); err != nil {
			return err
		}
		// TODO error
		// if _, err := conn.Write([]byte("unknown command")); err != nil {
		// 	return err
		// }
	}

	return nil
}
