package fields

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestChecksum_NewsEqual(t *testing.T) {
	tests := []struct {
		c1 Checksum
		c2 Checksum
	}{
		{NewChecksum([]byte{0x00, 0x01, 0x02}), NewChecksumWithParts([]byte{0x00}, []byte{0x01, 0x02})},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c1),
			func(t *testing.T) {
				if tt.c1 != tt.c2 {
					t.Errorf("NewChecksum: %v, NewChecsumWithParts: %v", tt.c1, tt.c2)
				}
			},
		)
	}
}

func TestChecksum_Size(t *testing.T) {
	tests := []struct {
		c    Checksum
		want int
	}{
		{Checksum(0), 4},
		{Checksum(3), 4},
		{Checksum(math.MaxUint32), 4},
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

func TestChecksum_Encode(t *testing.T) {
	tests := []struct {
		c    Checksum
		want []byte
	}{
		{Checksum(0), []byte{0x00, 0x00, 0x00, 0x00}},
		{Checksum(3), []byte{0x00, 0x00, 0x00, 0x03}},
		{Checksum(math.MaxUint32), []byte{0xFF, 0xFF, 0xFF, 0xFF}},
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
