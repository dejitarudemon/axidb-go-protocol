package values

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestInt(t *testing.T) {
	tests := []struct {
		name    string
		v       Int
		want    []byte
		wantErr bool
	}{
		{"zero", Int(0), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, false},
		{"positive", Int(10), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0A}, false},
		{"negative", Int(-10), []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF6}, false},
		{"max", Int(math.MaxInt64), []byte{0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"min", Int(math.MinInt64), []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValue(t, tt.v, fields.Int, tt.want, tt.wantErr)
		})
	}
}
