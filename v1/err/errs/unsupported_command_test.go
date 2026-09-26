package errs

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorUnsupportedCommand(t *testing.T) {
	tests := []struct {
		name    string
		got     fields.Command
		id      fields.TracebackID
		want    []byte
		wantErr bool
	}{
		{"reserved command", fields.Command(7), tracebackIDOne, encoded(3, tracebackIDOne), false},
		{"another reserved command", fields.Command(8), tracebackIDOne, encoded(3, tracebackIDOne), false},
		{"max command", fields.Command(255), tracebackIDMax, encoded(3, tracebackIDMax), false},
		{"handshake", fields.Handshake, tracebackIDOne, encoded(3, tracebackIDOne), true},
		{"answer", fields.Answer, tracebackIDOne, encoded(3, tracebackIDOne), true},
		{"read", fields.Read, tracebackIDOne, encoded(3, tracebackIDOne), true},
		{"write", fields.Write, tracebackIDOne, encoded(3, tracebackIDOne), true},
		{"delete", fields.Delete, tracebackIDOne, encoded(3, tracebackIDOne), true},
		{"batch", fields.Batch, tracebackIDOne, encoded(3, tracebackIDOne), true},
		{"ping", fields.Ping, tracebackIDOne, encoded(3, tracebackIDOne), true},
		{"zero value", fields.Handshake, fields.TracebackID{}, encoded(3, fields.TracebackID{}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorUnsupportedCommandWithTracebackID(tt.got, tt.id)
			assertProtocolError(t, e, fields.UnsupportedCommand, tt.id, tt.want, tt.wantErr)
		})
	}
}
