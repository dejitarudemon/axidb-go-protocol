package headers

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compression"
)

func TestHeaders_Size(t *testing.T) {
	tests := []struct {
		h    Headers
		want int
	}{
		{Headers{}, 10},
		{Headers{Compression: compression.Lz4}, 10},
		{Headers{Command: command.Read}, 10},
		{Headers{BodyLen: 1}, 10},
		{Headers{RequestID: 2}, 10},
		{Headers{Compression: compression.Lz4, Command: command.Read, BodyLen: 2, RequestID: 1}, 10},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestHeaders_Size %v", tt.h),
			func(t *testing.T) {
				if got := tt.h.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestHeaders_Encode(t *testing.T) {
	tests := []struct {
		h    Headers
		want []byte
	}{
		{
			Headers{},
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			Headers{Compression: compression.Lz4},
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00},
		},
		{
			Headers{Command: command.Read},
			[]byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			Headers{BodyLen: 1},
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			Headers{RequestID: 2},
			[]byte{0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			Headers{Compression: compression.Lz4, Command: command.Read, BodyLen: 2, RequestID: 1},
			[]byte{0x02, 0x00, 0x00, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x02},
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestHeaders_Encode %v", tt.h),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.h.Size())

				tt.h.Encode(&buf)

				got := buf.Bytes()

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestHeaders_IsValid(t *testing.T) {
	tests := []struct {
		h    Headers
		want bool
	}{
		{Headers{}, false},
		{Headers{Compression: compression.Lz4}, false},
		{Headers{Command: command.Read}, false},
		{Headers{BodyLen: 1}, false},
		{Headers{RequestID: 2}, false},
		{Headers{Compression: compression.Lz4, Command: command.Read, BodyLen: 2, RequestID: 1}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestHeaders_IsValid %v", tt.h),
			func(t *testing.T) {
				if got := tt.h.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
