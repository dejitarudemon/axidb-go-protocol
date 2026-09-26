package fields

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func TestBatchLimit(t *testing.T) {
	tests := []struct {
		name string
		b    BatchLimit
		want []byte
	}{
		{"zero", BatchLimit(0), []byte{0x00, 0x00, 0x00, 0x00}},
		{"small", BatchLimit(3), []byte{0x00, 0x00, 0x00, 0x03}},
		{"max", BatchLimit(math.MaxUint32), []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.AssertEncoded(t, tt.b, tt.want)
		})
	}
}
