package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

type MockConn struct {
	ResultOfRead  int
	ResultOfWrite int
	Buf           *bytes.Buffer
	ErrRead       error
	ErrWrite      error
	ErrClose      error
}

func (c *MockConn) Read(p []byte) (n int, err error) {
	return c.ResultOfRead, c.ErrRead
}
func (c *MockConn) Write(p []byte) (n int, err error) {
	if c.Buf == nil {
		c.Buf = &bytes.Buffer{}
	}
	if _, err := c.Buf.Write(p); err != nil {
		panic(fmt.Sprintf("failed to write buffer. err: %v", err))
	}
	return c.ResultOfWrite, c.ErrWrite
}
func (c *MockConn) Close() error {
	return c.ErrClose
}
func (c *MockConn) GetBuf() []byte {
	return c.Buf.Bytes()
}

func Test_handler(t *testing.T) {
	type args struct {
		conn io.ReadWriteCloser
	}
	tests := []struct {
		name     string
		args     args
		expected []byte
		wantErr  bool
	}{
		{
			name: "returns PONG for PING",
			args: args{
				conn: &MockConn{},
			},
			expected: []byte("+PONG\r\n"),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handler(tt.args.conn); (err != nil) != tt.wantErr {
				t.Errorf("handler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
