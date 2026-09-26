package fields

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func TestRequestID(t *testing.T) {
	tests := []struct {
		name string
		r    RequestID
		want []byte
	}{
		{"zero", RequestID(0), []byte{0x00, 0x00, 0x00, 0x00}},
		{"small", RequestID(3), []byte{0x00, 0x00, 0x00, 0x03}},
		{"max", RequestID(math.MaxUint32), []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.AssertEncoded(t, tt.r, tt.want)
		})
	}
}
