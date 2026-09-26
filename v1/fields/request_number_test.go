package fields

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func TestRequestNumber(t *testing.T) {
	tests := []struct {
		name string
		r    RequestNumber
		want []byte
	}{
		{"zero", RequestNumber(0), []byte{0x00, 0x00, 0x00, 0x00}},
		{"small", RequestNumber(3), []byte{0x00, 0x00, 0x00, 0x03}},
		{"max", RequestNumber(math.MaxUint32), []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.AssertEncoded(t, tt.r, tt.want)
		})
	}
}
