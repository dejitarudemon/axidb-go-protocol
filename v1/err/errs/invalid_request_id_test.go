package errs

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorInvalidRequestID(t *testing.T) {
	tests := []struct {
		name      string
		command   fields.Command
		requestID fields.RequestID
		id        fields.TracebackID
		want      []byte
		wantErr   bool
	}{
		{"answer without request id", fields.Answer, 0, tracebackIDOne, encoded(15, tracebackIDOne), false},
		{"answer with request id", fields.Answer, 1, tracebackIDOne, encoded(15, tracebackIDOne), false},
		{"batch without request id", fields.Batch, 0, tracebackIDOne, encoded(15, tracebackIDOne), false},
		{"batch with request id", fields.Batch, 1, tracebackIDOne, encoded(15, tracebackIDOne), true},
		{"handshake with max request id", fields.Handshake, math.MaxUint32, tracebackIDMax, encoded(15, tracebackIDMax), false},
		{"custom command with request id", fields.Command(10), 1, tracebackIDOne, encoded(15, tracebackIDOne), true},
		{"custom command with max request id", fields.Command(11), math.MaxUint32, tracebackIDOne, encoded(15, tracebackIDOne), true},
		{"max command with max request id", fields.Command(255), math.MaxUint32, tracebackIDOne, encoded(15, tracebackIDOne), true},
		{"zero value", fields.Handshake, 0, fields.TracebackID{}, encoded(15, fields.TracebackID{}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorInvalidRequestIDWithTracebackID(tt.command, tt.requestID, tt.id)
			assertProtocolError(t, e, fields.InvalidRequestID, tt.id, tt.want, tt.wantErr)
		})
	}
}
