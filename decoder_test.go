package main

import (
	"io"
	"reflect"
	"testing"
)

func TestRESPDecoder_Decode(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		d       *RESPDecoder
		args    args
		want    *DecodedMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &RESPDecoder{}
			got, err := d.Decode(tt.args.reader)
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
