package fields

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestType_Size(t *testing.T) {
	tests := []struct {
		t    Type
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
		{Type(7), 1},
		{Type(255), 1},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.t),
			func(t *testing.T) {
				if got := tt.t.Size(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestType_Encode(t *testing.T) {
	tests := []struct {
		t    Type
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
		{Type(8), []byte{0x08}},
		{Type(255), []byte{0xFF}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.t),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.t.Size())

				tt.t.Encode(&buf)
				got := buf.Bytes()

				if len(got) != tt.t.Size() {
					t.Fatalf("expected %v bytes, got %v bytes", tt.t.Size(), len(got))
				}

				if !bytes.Equal(got, tt.want) {
					t.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestType_String(t *testing.T) {
	tests := []struct {
		t    Type
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
		{Type(8), "Unknown (8)"},
		{Type(255), "Unknown (255)"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.t),
			func(t *testing.T) {
				if got := tt.t.String(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestType_IsValid(t *testing.T) {
	tests := []struct {
		t    Type
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
		{Type(8), false},
		{Type(255), false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.t),
			func(t *testing.T) {
				if got := tt.t.IsValid(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
