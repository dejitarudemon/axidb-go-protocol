package types

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestCode_Size(t *testing.T) {
	tests := []struct {
		c    Code
		want int
	}{
		{Bytes, 1},
		{TypedArray, 1},
		{UntypedArray, 1},
		{Int, 1},
		{Uint, 1},
		{Float, 1},
		{String, 1},
		{JSON, 1},
		{Code(7), 1},
		{Code(255), 1},
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

func TestCode_Encode(t *testing.T) {
	tests := []struct {
		c    Code
		want []byte
	}{
		{Bytes, []byte{0x00}},
		{TypedArray, []byte{0x01}},
		{UntypedArray, []byte{0x02}},
		{Int, []byte{0x03}},
		{Uint, []byte{0x04}},
		{Float, []byte{0x05}},
		{String, []byte{0x06}},
		{JSON, []byte{0x07}},
		{Code(8), []byte{0x08}},
		{Code(255), []byte{0xFF}},
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

func TestCode_String(t *testing.T) {
	tests := []struct {
		c    Code
		want string
	}{
		{Bytes, "Bytes"},
		{TypedArray, "Typed Array"},
		{UntypedArray, "Untyped Array"},
		{Int, "Int"},
		{Uint, "Uint"},
		{Float, "Float"},
		{String, "String"},
		{JSON, "JSON"},
		{Code(8), "Unknown (8)"},
		{Code(255), "Unknown (255)"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.String(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestCode_IsValid(t *testing.T) {
	tests := []struct {
		c    Code
		want bool
	}{
		{Bytes, true},
		{TypedArray, true},
		{UntypedArray, true},
		{Int, true},
		{Uint, true},
		{Float, true},
		{String, true},
		{JSON, true},
		{Code(8), false},
		{Code(255), false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.IsValid(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
