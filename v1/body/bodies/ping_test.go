package bodies

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name    string
		b       Ping
		want    []byte
		wantErr bool
	}{
		{"empty", Ping{}, []byte{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBody(t, tt.b, fields.Ping, tt.want, tt.wantErr)
		})
	}
}
