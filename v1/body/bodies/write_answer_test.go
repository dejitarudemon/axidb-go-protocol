package bodies

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestWriteAnswer(t *testing.T) {
	tests := []struct {
		name    string
		b       WriteAnswer
		want    []byte
		wantErr bool
	}{
		{"empty", WriteAnswer{}, []byte{0x01, 0x03}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertAnswer(t, tt.b, fields.Write, tt.want, tt.wantErr)
		})
	}
}
