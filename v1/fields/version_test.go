package fields

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestVersion_Size(t *testing.T) {
	tests := []struct {
		r    Version
		want int
	}{
		{Version(0), 1},
		{Version(3), 1},
		{Version(255), 1},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.r),
			func(t *testing.T) {
				if got := tt.r.Size(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestVersion_Encode(t *testing.T) {
	tests := []struct {
		r    Version
		want []byte
	}{
		{Version(0), []byte{0x00}},
		{Version(3), []byte{0x03}},
		{Version(math.MaxUint32), []byte{0xFF}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.r),
			func(t *testing.T) {
				buf := buffer.Slice{}
				buf.Preallocate(tt.r.Size())

				tt.r.Encode(&buf)
				got := buf.Bytes()

				if len(got) != tt.r.Size() {
					t.Fatalf("expected %v bytes, got %v bytes", tt.r.Size(), len(got))
				}

				if !bytes.Equal(got, tt.want) {
					t.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			},
		)
	}
}
