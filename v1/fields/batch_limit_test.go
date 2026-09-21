package fields

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestBatchLimit_Size(t *testing.T) {
	tests := []struct {
		b    BatchLimit
		want int
	}{
		{BatchLimit(0), 4},
		{BatchLimit(3), 4},
		{BatchLimit(math.MaxUint32), 4},
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

func TestBatchLimit_Encode(t *testing.T) {
	tests := []struct {
		b    BatchLimit
		want []byte
	}{
		{BatchLimit(0), []byte{0x00, 0x00, 0x00, 0x00}},
		{BatchLimit(3), []byte{0x00, 0x00, 0x00, 0x03}},
		{BatchLimit(math.MaxUint32), []byte{0xFF, 0xFF, 0xFF, 0xFF}},
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
