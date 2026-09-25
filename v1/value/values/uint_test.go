package values

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestUint_Size(t *testing.T) {
	tests := []struct {
		c    Uint
		want int
	}{
		{Uint(0), 8},
		{Uint(math.MaxUint64), 8},
		{Uint(10), 8},
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

func TestUint_Encode(t *testing.T) {
	tests := []struct {
		c    Uint
		want []byte
	}{
		{Uint(0), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{Uint(math.MaxUint64), []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{Uint(10), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0A}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				buf := buffer.Slice{}
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

func TestUint_Type(t *testing.T) {
	tests := []struct {
		c    Uint
		want fields.Type
	}{
		{Uint(0), fields.Uint},
		{Uint(math.MaxUint64), fields.Uint},
		{Uint(10), fields.Uint},
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

func TestUint_IsValid(t *testing.T) {
	tests := []struct {
		c       Uint
		wantErr bool
	}{
		{Uint(0), false},
		{Uint(math.MaxUint64), false},
		{Uint(10), false},
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
