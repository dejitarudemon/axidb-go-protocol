package values

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestBytes_Size(t *testing.T) {
	tests := []struct {
		c    Bytes
		want int
	}{
		{Bytes([]byte("")), 4},
		{Bytes([]byte("a")), 5},
		{Bytes([]byte{}), 4},
		{Bytes(nil), 4},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.Size(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestBytes_Encode(t *testing.T) {
	tests := []struct {
		c    Bytes
		want []byte
	}{
		{Bytes([]byte("")), []byte{0x00, 0x00, 0x00, 0x00}},
		{Bytes([]byte{0x0a, 0xff}), []byte{0x00, 0x00, 0x00, 0x02, 0x0a, 0xff}},
		{Bytes([]byte{}), []byte{0x00, 0x00, 0x00, 0x00}},
		{Bytes(nil), []byte{0x00, 0x00, 0x00, 0x00}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.c.Size())

				tt.c.Encode(&buf)
				got := buf.Bytes()

				if len(got) != tt.c.Size() {
					t.Fatalf("expected %v bytes, got %v bytes", tt.c.Size(), len(got))
				}

				if !bytes.Equal(got, tt.want) {
					t.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestBytes_Type(t *testing.T) {
	tests := []struct {
		c    Bytes
		want fields.Type
	}{
		{Bytes([]byte("")), fields.Bytes},
		{Bytes([]byte{0x0a, 0xff}), fields.Bytes},
		{Bytes([]byte{}), fields.Bytes},
		{Bytes(nil), fields.Bytes},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				got := tt.c.Type()

				if got != tt.c.Type() {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestBytes_IsValid(t *testing.T) {
	tests := []struct {
		c       Bytes
		wantErr bool
	}{
		{Bytes([]byte("")), false},
		{Bytes([]byte{0x0a, 0xff}), false},
		{Bytes([]byte{}), false},
		{Bytes(nil), false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.IsValid(); got == nil == tt.wantErr {
					t.Errorf("got %v, want %v", got, tt.wantErr)
				}
			},
		)
	}
}
