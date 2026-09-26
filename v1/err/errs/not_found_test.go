package errs

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorNotFound(t *testing.T) {
	tests := []struct {
		name    string
		key     fields.Key
		id      fields.TracebackID
		want    []byte
		wantErr bool
	}{
		{"nil key", nil, tracebackIDOne, encoded(10, tracebackIDOne), true},
		{"empty key", fields.Key(""), tracebackIDOne, encoded(10, tracebackIDOne), true},
		{"key", fields.Key("key"), tracebackIDMax, encoded(10, tracebackIDMax), false},
		{"zero value", nil, fields.TracebackID{}, encoded(10, fields.TracebackID{}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorNotFoundWithTracebackID(tt.key, tt.id)
			assertProtocolError(t, e, fields.NotFound, tt.id, tt.want, tt.wantErr)
		})
	}
}
