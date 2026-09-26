package errs

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorRequestsConflict(t *testing.T) {
	tests := []struct {
		name    string
		got     fields.RequestID
		id      fields.TracebackID
		want    []byte
		wantErr bool
	}{
		{"zero request id", 0, tracebackIDOne, encoded(4, tracebackIDOne), false},
		{"request id", 1, tracebackIDOne, encoded(4, tracebackIDOne), false},
		{"max request id", math.MaxUint32, tracebackIDMax, encoded(4, tracebackIDMax), false},
		{"zero value", 0, fields.TracebackID{}, encoded(4, fields.TracebackID{}), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorRequestsConflictWithTracebackID(tt.got, tt.id)
			assertProtocolError(t, e, fields.RequestsConflict, tt.id, tt.want, tt.wantErr)
		})
	}
}
