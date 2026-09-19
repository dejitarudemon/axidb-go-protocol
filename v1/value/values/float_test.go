package values

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/types"
)

func TestFloat_Size(t *testing.T) {
	tests := []struct {
		c    Float
		want int
	}{
		{Float(0), 8},
		{Float(math.MaxFloat64), 8},
		{Float(math.SmallestNonzeroFloat64), 8},
		{Float(-math.MaxFloat64), 8},
		{Float(-math.SmallestNonzeroFloat64), 8},
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

func TestFloat_Encode(t *testing.T) {
	tests := []struct {
		c    Float
		want []byte
	}{
		{Float(0), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{Float(math.MaxFloat64), []byte{0x7F, 0xEF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{Float(math.SmallestNonzeroFloat64), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}},
		{Float(-math.MaxFloat64), []byte{0xFF, 0xEF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{Float(-math.SmallestNonzeroFloat64), []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}},
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

func TestFloat_Type(t *testing.T) {
	tests := []struct {
		c    Float
		want types.Code
	}{
		{Float(0), types.Float},
		{Float(math.MaxFloat64), types.Float},
		{Float(math.SmallestNonzeroFloat64), types.Float},
		{Float(-math.MaxFloat64), types.Float},
		{Float(-math.SmallestNonzeroFloat64), types.Float},
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

func TestFloat_IsValid(t *testing.T) {
	tests := []struct {
		c       Float
		wantErr bool
	}{
		{Float(0), false},
		{Float(math.MaxFloat64), false},
		{Float(math.SmallestNonzeroFloat64), false},
		{Float(-math.MaxFloat64), false},
		{Float(-math.SmallestNonzeroFloat64), false},
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
