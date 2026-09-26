package errs

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorNoHello(t *testing.T) {
	tests := []struct {
		name    string
		id      fields.TracebackID
		want    []byte
		wantErr bool
	}{
		{"zero traceback id", fields.TracebackID{}, encoded(0, fields.TracebackID{}), false},
		{"traceback id", tracebackIDOne, encoded(0, tracebackIDOne), false},
		{"max traceback id", tracebackIDMax, encoded(0, tracebackIDMax), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorNoHelloWithTracebackID(tt.id)
			assertProtocolError(t, e, fields.NoHello, tt.id, tt.want, tt.wantErr)
		})
	}
}
