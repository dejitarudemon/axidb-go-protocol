package bodies

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestPingAnswer(t *testing.T) {
	tests := []struct {
		name    string
		b       PingAnswer
		want    []byte
		wantErr bool
	}{
		{"empty", PingAnswer{}, []byte{0x01, 0x06}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertAnswer(t, tt.b, fields.Ping, tt.want, tt.wantErr)
		})
	}
}
