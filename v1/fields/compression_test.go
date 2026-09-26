package fields

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func TestCompression(t *testing.T) {
	tests := []struct {
		c         Compression
		want      []byte
		wantStr   string
		wantValid bool
	}{
		{None, []byte{0x00}, "None", true},
		{Zstd, []byte{0x01}, "Zstd", true},
		{S2, []byte{0x02}, "S2", true},
		{Compression(3), []byte{0x03}, "Unknown (3)", false},
		{Compression(255), []byte{0xFF}, "Unknown (255)", false},
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
