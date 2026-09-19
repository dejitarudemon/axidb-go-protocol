package values

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/types"
)

func TestInt_Size(t *testing.T) {
	tests := []struct {
		c    Int
		want int
	}{
		{Int(0), 8},
		{Int(math.MaxInt64), 8},
		{Int(10), 8},
		{Int(math.MinInt64), 8},
		{Int(-10), 8},
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

func TestInt_Encode(t *testing.T) {
	tests := []struct {
		c    Int
		want []byte
	}{
		{Int(0), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{Int(math.MaxInt64), []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{Int(10), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0A}},
		{Int(math.MinInt64), []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{Int(-10), []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF6}},
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

func TestInt_Type(t *testing.T) {
	tests := []struct {
		c    Int
		want types.Code
	}{
		{Int(0), types.Int},
		{Int(math.MaxInt64), types.Int},
		{Int(10), types.Int},
		{Int(math.MinInt), types.Int},
		{Int(-10), types.Int},
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

func TestInt_IsValid(t *testing.T) {
	tests := []struct {
		c       Int
		wantErr bool
	}{
		{Int(0), false},
		{Int(math.MaxInt64), false},
		{Int(10), false},
		{Int(math.MinInt64), false},
		{Int(-10), false},
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
