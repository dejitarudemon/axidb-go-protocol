package values

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestString(t *testing.T) {
	tests := []struct {
		name    string
		v       String
		want    []byte
		wantErr bool
	}{
		{"empty", String(""), []byte{0x00, 0x00, 0x00, 0x00}, false},
		{"ascii", String("a"), []byte{0x00, 0x00, 0x00, 0x01, 0x61}, false},
		{"multibyte", String("ф"), []byte{0x00, 0x00, 0x00, 0x02, 0xD1, 0x84}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValue(t, tt.v, fields.String, tt.want, tt.wantErr)
		})
	}
}
