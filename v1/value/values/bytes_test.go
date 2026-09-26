package values

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestBytes(t *testing.T) {
	tests := []struct {
		name    string
		v       Bytes
		want    []byte
		wantErr bool
	}{
		{"nil", Bytes(nil), []byte{0x00, 0x00, 0x00, 0x00}, false},
		{"empty", Bytes([]byte{}), []byte{0x00, 0x00, 0x00, 0x00}, false},
		{"one byte", Bytes("a"), []byte{0x00, 0x00, 0x00, 0x01, 0x61}, false},
		{"two bytes", Bytes([]byte{0x0A, 0xFF}), []byte{0x00, 0x00, 0x00, 0x02, 0x0A, 0xFF}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValue(t, tt.v, fields.Bytes, tt.want, tt.wantErr)
		})
	}
}
