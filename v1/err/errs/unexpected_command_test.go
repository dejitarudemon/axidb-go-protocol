package errs

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorUnexpectedCommand(t *testing.T) {
	tests := []struct {
		name     string
		got      fields.Command
		expected fields.Command
		id       fields.TracebackID
		want     []byte
		wantErr  bool
	}{
		{"answer instead of handshake", fields.Answer, fields.Handshake, tracebackIDOne, encoded(2, tracebackIDOne, 0x00), false},
		{"handshake instead of read", fields.Handshake, fields.Read, tracebackIDMax, encoded(2, tracebackIDMax, 0x02), false},
		{"read instead of write", fields.Read, fields.Write, tracebackIDOne, encoded(2, tracebackIDOne, 0x03), false},
		{"write instead of delete", fields.Write, fields.Delete, tracebackIDOne, encoded(2, tracebackIDOne, 0x04), false},
		{"delete instead of ping", fields.Delete, fields.Ping, tracebackIDOne, encoded(2, tracebackIDOne, 0x06), false},
		{"ping instead of batch", fields.Ping, fields.Batch, tracebackIDOne, encoded(2, tracebackIDOne, 0x05), false},
		{"batch instead of answer", fields.Batch, fields.Answer, tracebackIDOne, encoded(2, tracebackIDOne, 0x01), false},
		{"reserved commands", fields.Command(7), fields.Command(10), tracebackIDOne, encoded(2, tracebackIDOne, 0x0A), false},
		{"max command", fields.Command(255), fields.Answer, tracebackIDOne, encoded(2, tracebackIDOne, 0x01), false},
		{"same commands", fields.Answer, fields.Answer, tracebackIDOne, encoded(2, tracebackIDOne, 0x01), true},
		{"zero value", fields.Handshake, fields.Handshake, fields.TracebackID{}, encoded(2, fields.TracebackID{}, 0x00), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorUnexpectedCommandWithTracebackID(tt.got, tt.expected, tt.id)
			assertProtocolError(t, e, fields.UnexpectedCommand, tt.id, tt.want, tt.wantErr)
		})
	}
}
