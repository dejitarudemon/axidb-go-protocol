package bodies

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		b       Read
		want    []byte
		wantErr bool
	}{
		{"empty key", Read{}, []byte{}, true},
		{"zero byte key", Read{0x00}, []byte{0x00}, false},
		{"ascii key", Read("ab"), []byte{0x61, 0x62}, false},
		{"utf-8 key", Read("фи"), []byte{0xD1, 0x84, 0xD0, 0xB8}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBody(t, tt.b, fields.Read, tt.want, tt.wantErr)
		})
	}
}
