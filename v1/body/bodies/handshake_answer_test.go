package bodies

import (
	"bytes"
	"fmt"
	"slices"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compression"
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
			NewHandshakeAnswer([]compression.Code{0x01, 0x00, 0x02}),
			HandshakeAnswer{[]compression.Code{0x01, 0x00, 0x02}},
		},
		{
			NewHandshakeAnswer([]compression.Code{0x01, 0x00, 0x02, 0x01, 0x03}),
			HandshakeAnswer{[]compression.Code{0x01, 0x00, 0x02, 0x03}},
		},
		{
			NewHandshakeAnswer([]compression.Code{0x01, 0x01}),
			HandshakeAnswer{[]compression.Code{0x01}},
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestHandshakeAnswer_New %v", tt.h),
			func(t *testing.T) {
				if !slices.Equal(tt.h.Compressions, tt.want.Compressions) {
					t.Errorf("compression: got '%v', want '%v'", tt.h.Compressions, tt.want.Compressions)
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
		{HandshakeAnswer{[]compression.Code{}}, 2},
		{HandshakeAnswer{[]compression.Code{0x00}}, 3},
		{HandshakeAnswer{[]compression.Code{0x00, 0x01}}, 4},
		{HandshakeAnswer{[]compression.Code{0x00, 0x01, 0x00}}, 5},
		{NewHandshakeAnswer([]compression.Code{0x00, 0x01, 0x00}), 4},
		{HandshakeAnswer{generateManyCompressions(1000)}, 258},
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
		{HandshakeAnswer{[]compression.Code{}}, []byte{0x01, 0x00}},
		{HandshakeAnswer{[]compression.Code{0x00}}, []byte{0x01, 0x01, 0x00}},
		{HandshakeAnswer{[]compression.Code{0x00, 0x01}}, []byte{0x01, 0x02, 0x00, 0x01}},
		{HandshakeAnswer{[]compression.Code{0x00, 0x01, 0x00}}, []byte{0x01, 0x03, 0x00, 0x01, 0x00}},
		{NewHandshakeAnswer([]compression.Code{0x00, 0x01, 0x00}), []byte{0x01, 0x02, 0x00, 0x01}},
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
		want command.Code
	}{
		{HandshakeAnswer{}, command.Answer},
		{HandshakeAnswer{[]compression.Code{}}, command.Answer},
		{HandshakeAnswer{[]compression.Code{0x00}}, command.Answer},
		{HandshakeAnswer{[]compression.Code{0x00, 0x01}}, command.Answer},
		{HandshakeAnswer{[]compression.Code{0x00, 0x01, 0x00}}, command.Answer},
		{NewHandshakeAnswer([]compression.Code{0x00, 0x01, 0x00}), command.Answer},
		{HandshakeAnswer{generateManyCompressions(1000)}, command.Answer},
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
		{HandshakeAnswer{[]compression.Code{}}, false},
		{HandshakeAnswer{[]compression.Code{0x00}}, false},
		{HandshakeAnswer{[]compression.Code{0x00, 0x01}}, false},
		{HandshakeAnswer{[]compression.Code{0x00, 0x01, 0x00}}, false},
		{NewHandshakeAnswer([]compression.Code{0x00, 0x01, 0x00}), false},
		{HandshakeAnswer{generateManyCompressions(1000)}, true},
		{NewHandshakeAnswer(generateManyCompressions(1000)), false},
		{HandshakeAnswer{[]compression.Code{0x00, 0x01, 0xFF}}, false},
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
