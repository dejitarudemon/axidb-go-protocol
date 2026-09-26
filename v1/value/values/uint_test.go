package values

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestUint(t *testing.T) {
	tests := []struct {
		name    string
		v       Uint
		want    []byte
		wantErr bool
	}{
		{"zero", Uint(0), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, false},
		{"small", Uint(10), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0A}, false},
		{"max", Uint(math.MaxUint64), []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValue(t, tt.v, fields.Uint, tt.want, tt.wantErr)
		})
	}
}
