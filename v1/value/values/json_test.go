package values

import (
	"bytes"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestJSON(t *testing.T) {
	invalid := func(n int) []byte { return bytes.Repeat([]byte{0xFF}, n) }
	prefixed := func(n int) []byte { return append([]byte{0x00, 0x00, 0x00, byte(n)}, invalid(n)...) }

	tests := []struct {
		name    string
		v       JSON
		want    []byte
		wantErr bool
	}{
		{"empty object", JSON(`{}`), []byte{0x00, 0x00, 0x00, 0x02, 0x7B, 0x7D}, false},
		{"string field", JSON(`{"1":"2"}`), []byte{0x00, 0x00, 0x00, 0x09, 0x7B, 0x22, 0x31, 0x22, 0x3A, 0x22, 0x32, 0x22, 0x7D}, false},
		{"number field", JSON(`{"1":1}`), []byte{0x00, 0x00, 0x00, 0x07, 0x7B, 0x22, 0x31, 0x22, 0x3A, 0x31, 0x7D}, false},
		{"empty", JSON([]byte{}), []byte{0x00, 0x00, 0x00, 0x00}, true},
		{"truncated", JSON(`{"1":`), []byte{0x00, 0x00, 0x00, 0x05, 0x7B, 0x22, 0x31, 0x22, 0x3A}, true},
		{"binary", JSON([]byte{0x00, 0x01, 0x02}), []byte{0x00, 0x00, 0x00, 0x03, 0x00, 0x01, 0x02}, true},
		{"invalid at detail limit", JSON(invalid(FirstSymbolsToShowJSON)), prefixed(FirstSymbolsToShowJSON), true},
		{"invalid over detail limit", JSON(invalid(FirstSymbolsToShowJSON + 1)), prefixed(FirstSymbolsToShowJSON + 1), true},
		{"long invalid", JSON(invalid(64)), prefixed(64), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValue(t, tt.v, fields.JSON, tt.want, tt.wantErr)
		})
	}
}
