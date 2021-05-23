package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestRESPDecoder_Decode(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		want    []string
		wantErr bool
	}{
		{
			name:    "string",
			msg:     "+PING\r\n",
			want:    []string{"PING"},
			wantErr: false,
		},
		{
			name:    "integer",
			msg:     ":1\r\n",
			want:    []string{"1"},
			wantErr: false,
		},
		{
			name:    "bulk string",
			msg:     "$3\r\nFOO\r\n",
			want:    []string{"FOO"},
			wantErr: false,
		},
		{
			name:    "arrays with 1 value",
			msg:     "*1\r\n$4\r\nPING\r\n",
			want:    []string{"PING"},
			wantErr: false,
		},
		{
			name:    "arrays with 1 values",
			msg:     "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n",
			want:    []string{"ECHO", "hey"},
			wantErr: false,
		},
		{
			name:    "error when no content",
			msg:     "+",
			wantErr: true,
		},
		{
			name:    "error when no crlf",
			msg:     "+OK",
			wantErr: true,
		},
		{
			name:    "error when invalid crlf",
			msg:     "+OK\r",
			wantErr: true,
		},
		{
			name:    "error when invalid crlf",
			msg:     "+OK\n",
			wantErr: true,
		},
		{
			name:    "error when invalid size",
			msg:     "*1\r\n$3\r\nPING\r\n",
			wantErr: true,
		},
		{
			name:    "error when no content in array",
			msg:     "*1\r\n$4",
			wantErr: true,
		},
		{
			name:    "error when invalid array size",
			msg:     "*2\r\n$4\r\nECHO\r\n",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewRESPDecoder(strings.NewReader(tt.msg))
			if err != nil {
				if !tt.wantErr {
					t.Errorf("RESPDecoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				return
			}

			got, err := d.Decode()
			if (err != nil) != tt.wantErr {
				t.Errorf("RESPDecoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RESPDecoder.Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}
