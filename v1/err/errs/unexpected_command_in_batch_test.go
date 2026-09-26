package errs

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorUnexpectedCommandInBatch(t *testing.T) {
	tests := []struct {
		name          string
		command       fields.Command
		requestNumber fields.RequestNumber
		id            fields.TracebackID
		want          []byte
		wantErr       bool
	}{
		{"answer", fields.Answer, 0, tracebackIDOne, encoded(14, tracebackIDOne, 0x00, 0x00, 0x00, 0x00), false},
		{"answer at request 1", fields.Answer, 1, tracebackIDOne, encoded(14, tracebackIDOne, 0x00, 0x00, 0x00, 0x01), false},
		{"batch", fields.Batch, 0, tracebackIDOne, encoded(14, tracebackIDOne, 0x00, 0x00, 0x00, 0x00), false},
		{"batch at request 1", fields.Batch, 1, tracebackIDOne, encoded(14, tracebackIDOne, 0x00, 0x00, 0x00, 0x01), false},
		{"handshake at max request", fields.Handshake, math.MaxUint32, tracebackIDMax, encoded(14, tracebackIDMax, 0xFF, 0xFF, 0xFF, 0xFF), false},
		{"ping at max request", fields.Ping, math.MaxUint32, tracebackIDOne, encoded(14, tracebackIDOne, 0xFF, 0xFF, 0xFF, 0xFF), false},
		{"custom command", fields.Command(10), 1, tracebackIDOne, encoded(14, tracebackIDOne, 0x00, 0x00, 0x00, 0x01), false},
		{"custom command at max request", fields.Command(11), math.MaxUint32, tracebackIDOne, encoded(14, tracebackIDOne, 0xFF, 0xFF, 0xFF, 0xFF), false},
		{"max command at max request", fields.Command(255), math.MaxUint32, tracebackIDOne, encoded(14, tracebackIDOne, 0xFF, 0xFF, 0xFF, 0xFF), false},
		{"read", fields.Read, 0, tracebackIDOne, encoded(14, tracebackIDOne, 0x00, 0x00, 0x00, 0x00), true},
		{"read at request 1", fields.Read, 1, tracebackIDOne, encoded(14, tracebackIDOne, 0x00, 0x00, 0x00, 0x01), true},
		{"write", fields.Write, 0, tracebackIDOne, encoded(14, tracebackIDOne, 0x00, 0x00, 0x00, 0x00), true},
		{"delete", fields.Delete, 2, tracebackIDOne, encoded(14, tracebackIDOne, 0x00, 0x00, 0x00, 0x02), true},
		{"zero value", fields.Handshake, 0, fields.TracebackID{}, encoded(14, fields.TracebackID{}, 0x00, 0x00, 0x00, 0x00), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorUnexpectedCommandInBatchWithTracebackID(tt.command, tt.requestNumber, tt.id)
			assertProtocolError(t, e, fields.UnexpectedCommandInBatch, tt.id, tt.want, tt.wantErr)
		})
	}
}
