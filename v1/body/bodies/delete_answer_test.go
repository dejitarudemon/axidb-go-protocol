package bodies

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestDeleteAnswer(t *testing.T) {
	tests := []struct {
		name    string
		b       DeleteAnswer
		want    []byte
		wantErr bool
	}{
		{"empty", DeleteAnswer{}, []byte{0x01, 0x04}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertAnswer(t, tt.b, fields.Delete, tt.want, tt.wantErr)
		})
	}
}
