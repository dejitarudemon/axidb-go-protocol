package compression

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
		{None, 1},
		{Zstd, 1},
		{Lz4, 1},
		{Code(3), 1},
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
		{None, []byte{0x00}},
		{Zstd, []byte{0x01}},
		{Lz4, []byte{0x02}},
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
		{None, "None"},
		{Zstd, "Zstd"},
		{Lz4, "Lz4"},
		{Code(3), "Unknown (3)"},
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
		c       Code
		wantErr bool
	}{
		{None, true},
		{Zstd, true},
		{Lz4, true},
		{Code(3), false},
		{Code(255), false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.IsValid(); got != tt.wantErr {
					t.Errorf("got %v, want %v", got, tt.wantErr)
				}
			},
		)
	}
}
