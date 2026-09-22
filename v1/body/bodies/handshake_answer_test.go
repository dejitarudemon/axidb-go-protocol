package bodies

import (
	"bytes"
	"fmt"
	"slices"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestHandshakeAnswer_New(t *testing.T) {
	tests := []struct {
		h    HandshakeAnswer
		want HandshakeAnswer
	}{
		{
			NewHandshakeAnswer(nil),
			HandshakeAnswer{nil},
		},
		{
			NewHandshakeAnswer([]fields.Compression{0x01, 0x00, 0x02}),
			HandshakeAnswer{[]fields.Compression{0x01, 0x00, 0x02}},
		},
		{
			NewHandshakeAnswer([]fields.Compression{0x01, 0x00, 0x02, 0x01, 0x03}),
			HandshakeAnswer{[]fields.Compression{0x01, 0x00, 0x02, 0x03}},
		},
		{
			NewHandshakeAnswer([]fields.Compression{0x01, 0x01}),
			HandshakeAnswer{[]fields.Compression{0x01}},
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestHandshakeAnswer_New %v", tt.h),
			func(t *testing.T) {
				if !slices.Equal(tt.h.Compressions, tt.want.Compressions) {
					t.Errorf("fields: got '%v', want '%v'", tt.h.Compressions, tt.want.Compressions)
				}
			},
		)
	}
}

func TestHandshakeAnswer_Size(t *testing.T) {
	tests := []struct {
		h    HandshakeAnswer
		want int
	}{
		{HandshakeAnswer{}, 2},
		{HandshakeAnswer{[]fields.Compression{}}, 2},
		{HandshakeAnswer{[]fields.Compression{0x00}}, 3},
		{HandshakeAnswer{[]fields.Compression{0x00, 0x01}}, 4},
		{HandshakeAnswer{[]fields.Compression{0x00, 0x01, 0x00}}, 5},
		{NewHandshakeAnswer([]fields.Compression{0x00, 0x01, 0x00}), 4},
		{HandshakeAnswer{generateManyCompressions(1000)}, 257},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestHandshakeAnswer_Size %v", tt.h),
			func(t *testing.T) {
				if got := tt.h.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestHandshakeAnswer_Encode(t *testing.T) {
	generated := generateManyCompressions(1000)
	tests := []struct {
		h    HandshakeAnswer
		want []byte
	}{
		{HandshakeAnswer{}, []byte{0x01, 0x00}},
		{HandshakeAnswer{[]fields.Compression{}}, []byte{0x01, 0x00}},
		{HandshakeAnswer{[]fields.Compression{0x00}}, []byte{0x01, 0x01, 0x00}},
		{HandshakeAnswer{[]fields.Compression{0x00, 0x01}}, []byte{0x01, 0x02, 0x00, 0x01}},
		{HandshakeAnswer{[]fields.Compression{0x00, 0x01, 0x00}}, []byte{0x01, 0x03, 0x00, 0x01, 0x00}},
		{NewHandshakeAnswer([]fields.Compression{0x00, 0x01, 0x00}), []byte{0x01, 0x02, 0x00, 0x01}},
		{HandshakeAnswer{generated}, append([]byte{0x01}, encodeCompressions(generated)...)},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestHandshakeAnswer_Encode %v", tt.h),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.h.Size())

				tt.h.Encode(&buf)

				got := buf.Bytes()

				if len(got) != len(tt.want) {
					t.Fatalf("got %v want %v", len(got), len(tt.want))
				}

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestHandshakeAnswer_Command(t *testing.T) {
	tests := []struct {
		h    HandshakeAnswer
		want fields.Command
	}{
		{HandshakeAnswer{}, fields.Answer},
		{HandshakeAnswer{[]fields.Compression{}}, fields.Answer},
		{HandshakeAnswer{[]fields.Compression{0x00}}, fields.Answer},
		{HandshakeAnswer{[]fields.Compression{0x00, 0x01}}, fields.Answer},
		{HandshakeAnswer{[]fields.Compression{0x00, 0x01, 0x00}}, fields.Answer},
		{NewHandshakeAnswer([]fields.Compression{0x00, 0x01, 0x00}), fields.Answer},
		{HandshakeAnswer{generateManyCompressions(1000)}, fields.Answer},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestHandshakeAnswer_Command %v", tt.h),
			func(t *testing.T) {
				if got := tt.h.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestHandshakeAnswer_IsValid(t *testing.T) {
	tests := []struct {
		h    HandshakeAnswer
		want bool
	}{
		{HandshakeAnswer{}, false},
		{HandshakeAnswer{[]fields.Compression{}}, false},
		{HandshakeAnswer{[]fields.Compression{0x00}}, false},
		{HandshakeAnswer{[]fields.Compression{0x00, 0x01}}, false},
		{HandshakeAnswer{[]fields.Compression{0x00, 0x01, 0x00}}, false},
		{NewHandshakeAnswer([]fields.Compression{0x00, 0x01, 0x00}), false},
		{HandshakeAnswer{generateManyCompressions(1000)}, true},
		{NewHandshakeAnswer(generateManyCompressions(1000)), false},
		{HandshakeAnswer{[]fields.Compression{0x00, 0x01, 0xFF}}, false},
	}
	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestHandshakeAnswer_IsValid %v", tt.h),
			func(t *testing.T) {
				if got := tt.h.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
