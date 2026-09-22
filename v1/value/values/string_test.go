package values

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestString_Size(t *testing.T) {
	tests := []struct {
		c    String
		want int
	}{
		{String(""), 4},
		{String("a"), 5},
		{String("ф"), 6},
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

func TestString_Encode(t *testing.T) {
	tests := []struct {
		c    String
		want []byte
	}{
		{String(""), []byte{0x00, 0x00, 0x00, 0x00}},
		{String("a"), []byte{0x00, 0x00, 0x00, 0x01, 0x61}},
		{String("ф"), []byte{0x00, 0x00, 0x00, 0x02, 0xD1, 0x84}},
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

func TestString_Type(t *testing.T) {
	tests := []struct {
		c    String
		want fields.Type
	}{
		{String(""), fields.String},
		{String("a"), fields.String},
		{String("ф"), fields.String},
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

func TestString_IsValid(t *testing.T) {
	tests := []struct {
		c       String
		wantErr bool
	}{
		{String(""), false},
		{String("a"), false},
		{String("ф"), false},
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
