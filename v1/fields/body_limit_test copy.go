package fields

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestBodyLimit_Size(t *testing.T) {
	tests := []struct {
		b    BodyLimit
		want int
	}{
		{BodyLimit(0), 4},
		{BodyLimit(3), 4},
		{BodyLimit(math.MaxUint32), 4},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.b),
			func(t *testing.T) {
				if got := tt.b.Size(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestBodyLimit_Encode(t *testing.T) {
	tests := []struct {
		b    BodyLimit
		want []byte
	}{
		{BodyLimit(0), []byte{0x00, 0x00, 0x00, 0x00}},
		{BodyLimit(3), []byte{0x00, 0x00, 0x00, 0x03}},
		{BodyLimit(math.MaxUint32), []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.b),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.b.Size())

				tt.b.Encode(&buf)
				got := buf.Bytes()

				if len(got) != tt.b.Size() {
					t.Fatalf("expected %v bytes, got %v bytes", tt.b.Size(), len(got))
				}

				if !bytes.Equal(got, tt.want) {
					t.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			},
		)
	}
}
