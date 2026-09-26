package fields

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		c         Command
		want      []byte
		wantStr   string
		wantValid bool
	}{
		{Handshake, []byte{0x00}, "Handshake", true},
		{Answer, []byte{0x01}, "Answer", true},
		{Read, []byte{0x02}, "Read", true},
		{Write, []byte{0x03}, "Write", true},
		{Delete, []byte{0x04}, "Delete", true},
		{Batch, []byte{0x05}, "Batch", true},
		{Ping, []byte{0x06}, "Ping", true},
		{Command(7), []byte{0x07}, "Unknown (7)", false},
		{Command(255), []byte{0xFF}, "Unknown (255)", false},
	}

	for _, tt := range tests {
		t.Run(tt.wantStr, func(t *testing.T) {
			testutil.AssertEncoded(t, tt.c, tt.want)

			if got := tt.c.String(); got != tt.wantStr {
				t.Errorf("String() = %q, want %q", got, tt.wantStr)
			}

			if got := tt.c.IsValid(); got != tt.wantValid {
				t.Errorf("IsValid() = %v, want %v", got, tt.wantValid)
			}
		})
	}
}
