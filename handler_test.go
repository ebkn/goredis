package main

import (
	"bytes"
	"fmt"
	"testing"
)

type MockConn struct {
	Req []byte
	Res *bytes.Buffer

	ResultOfRead  int
	ResultOfWrite int
	ErrRead       error
	ErrWrite      error
	ErrClose      error
}

func (c *MockConn) Read(p []byte) (n int, err error) {
	buf := bytes.NewBuffer(p)
	if _, err := buf.Write(c.Req); err != nil {
		panic(fmt.Sprintf("failed to write buffer. err: %v", err))
	}
	return c.ResultOfRead, c.ErrRead
}
func (c *MockConn) Write(p []byte) (n int, err error) {
	if c.Res == nil {
		c.Res = &bytes.Buffer{}
	}
	if _, err := c.Res.Write(p); err != nil {
		panic(fmt.Sprintf("failed to write buffer. err: %v", err))
	}
	return c.ResultOfWrite, c.ErrWrite
}
func (c *MockConn) Close() error {
	return c.ErrClose
}
func (c *MockConn) GetRes() []byte {
	return c.Res.Bytes()
}

func Test_handler(t *testing.T) {
	type args struct {
		conn *MockConn
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
				conn: &MockConn{
					Req: []byte("+PING\r\n"),
				},
			},
			expected: []byte("+PONG\r\n"),
			wantErr:  false,
		},
		{
			name: "returns given string for ECHO",
			args: args{
				conn: &MockConn{
					Req: []byte("+ECHO Hello\r\n"),
				},
			},
			expected: []byte("Hello\r\n"),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handler(tt.args.conn); (err != nil) != tt.wantErr {
				t.Errorf("handler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if bytes.Compare(tt.args.conn.GetRes(), tt.expected) != 0 {
				t.Errorf("result expected=%v, but got=%v", string(tt.expected), string(tt.args.conn.GetRes()))
			}
		})
	}
}
