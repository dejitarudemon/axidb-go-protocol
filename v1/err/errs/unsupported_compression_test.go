package errs

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorUnsupportedCompression(t *testing.T) {
	tests := []struct {
		name        string
		compression fields.Compression
		id          fields.TracebackID
		want        []byte
		wantErr     bool
	}{
		{"s2", fields.S2, tracebackIDOne, encoded(5, tracebackIDOne), false},
		{"zstd", fields.Zstd, tracebackIDOne, encoded(5, tracebackIDOne), false},
		{"custom compression", fields.Compression(3), tracebackIDOne, encoded(5, tracebackIDOne), false},
		{"max compression", fields.Compression(255), tracebackIDMax, encoded(5, tracebackIDMax), false},
		{"no compression", fields.None, tracebackIDOne, encoded(5, tracebackIDOne), true},
		{"zero value", fields.None, fields.TracebackID{}, encoded(5, fields.TracebackID{}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorUnsupportedCompressionWithTracebackID(tt.compression, tt.id)
			assertProtocolError(t, e, fields.UnsupportedCompression, tt.id, tt.want, tt.wantErr)
		})
	}
}
