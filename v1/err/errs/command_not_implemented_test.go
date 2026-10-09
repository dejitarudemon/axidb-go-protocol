package errs

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorCommandNotImplemented(t *testing.T) {
	tests := []struct {
		name    string
		id      fields.TracebackID
		want    []byte
		wantErr bool
	}{
		{"zero traceback id", fields.TracebackID{}, encoded(1, fields.TracebackID{}), false},
		{"traceback id", tracebackIDOne, encoded(1, tracebackIDOne), false},
		{"max traceback id", tracebackIDMax, encoded(1, tracebackIDMax), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorCommandNotImplementedWithTracebackID(tt.id)
			assertProtocolError(t, e, fields.CommandNotImplemented, tt.id, tt.want, tt.wantErr)
		})
	}
}
