package fields

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestRequestID_Size(t *testing.T) {
	tests := []struct {
		r    RequestID
		want int
	}{
		{RequestID(0), 4},
		{RequestID(3), 4},
		{RequestID(math.MaxUint32), 4},
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

func TestRequestID_Encode(t *testing.T) {
	tests := []struct {
		r    RequestID
		want []byte
	}{
		{RequestID(0), []byte{0x00, 0x00, 0x00, 0x00}},
		{RequestID(3), []byte{0x00, 0x00, 0x00, 0x03}},
		{RequestID(math.MaxUint32), []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.r),
			func(t *testing.T) {
				buf := buffer.Mock{}
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
