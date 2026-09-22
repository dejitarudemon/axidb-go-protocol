package fields

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestCompression_Size(t *testing.T) {
	tests := []struct {
		c    Compression
		want int
	}{
		{None, 1},
		{Zstd, 1},
		{S2, 1},
		{Compression(3), 1},
		{Compression(255), 1},
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

func TestCompression_Encode(t *testing.T) {
	tests := []struct {
		c    Compression
		want []byte
	}{
		{None, []byte{0x00}},
		{Zstd, []byte{0x01}},
		{S2, []byte{0x02}},
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

func TestCompression_String(t *testing.T) {
	tests := []struct {
		c    Compression
		want string
	}{
		{None, "None"},
		{Zstd, "Zstd"},
		{S2, "S2"},
		{Compression(3), "Unknown (3)"},
		{Compression(255), "Unknown (255)"},
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

func TestCompression_IsValid(t *testing.T) {
	tests := []struct {
		c       Compression
		wantErr bool
	}{
		{None, true},
		{Zstd, true},
		{S2, true},
		{Compression(3), false},
		{Compression(255), false},
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
