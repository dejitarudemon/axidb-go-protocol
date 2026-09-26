package errs

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorRequestInterrupted(t *testing.T) {
	tests := []struct {
		name        string
		interrupted fields.RequestID
		id          fields.TracebackID
		want        []byte
		wantErr     bool
	}{
		{"zero request id", 0, tracebackIDOne, encoded(12, tracebackIDOne), false},
		{"request id", 1, tracebackIDOne, encoded(12, tracebackIDOne), false},
		{"max request id", math.MaxUint32, tracebackIDMax, encoded(12, tracebackIDMax), false},
		{"zero value", 0, fields.TracebackID{}, encoded(12, fields.TracebackID{}), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorRequestInterruptedWithTracebackID(tt.interrupted, tt.id)
			assertProtocolError(t, e, fields.RequestInterrupted, tt.id, tt.want, tt.wantErr)
		})
	}
}
