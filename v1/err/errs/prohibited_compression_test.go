package errs

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorProhibitedCompression(t *testing.T) {
	tests := []struct {
		name        string
		compression fields.Compression
		command     fields.Command
		id          fields.TracebackID
		want        []byte
		wantErr     bool
	}{
		{"s2 for handshake", fields.S2, fields.Handshake, tracebackIDOne, encoded(11, tracebackIDOne), false},
		{"zstd for ping", fields.Zstd, fields.Ping, tracebackIDMax, encoded(11, tracebackIDMax), false},
		{"zstd for answer", fields.Zstd, fields.Answer, tracebackIDOne, encoded(11, tracebackIDOne), false},
		{"max compression for handshake", fields.Compression(255), fields.Handshake, tracebackIDOne, encoded(11, tracebackIDOne), false},
		{"no compression for handshake", fields.None, fields.Handshake, tracebackIDOne, encoded(11, tracebackIDOne), true},
		{"s2 for read", fields.S2, fields.Read, tracebackIDOne, encoded(11, tracebackIDOne), true},
		{"zstd for write", fields.Zstd, fields.Write, tracebackIDOne, encoded(11, tracebackIDOne), true},
		{"custom compression for delete", fields.Compression(4), fields.Delete, tracebackIDOne, encoded(11, tracebackIDOne), true},
		{"no compression for custom command", fields.None, fields.Command(4), tracebackIDOne, encoded(11, tracebackIDOne), true},
		{"custom compression for custom command", fields.Compression(5), fields.Command(5), tracebackIDOne, encoded(11, tracebackIDOne), true},
		{"max compression for max command", fields.Compression(255), fields.Command(255), tracebackIDOne, encoded(11, tracebackIDOne), true},
		{"zero value", fields.None, fields.Handshake, fields.TracebackID{}, encoded(11, fields.TracebackID{}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorProhibitedCompressionWithTracebackID(tt.compression, tt.command, tt.id)
			assertProtocolError(t, e, fields.ProhibitedCompression, tt.id, tt.want, tt.wantErr)
		})
	}
}
