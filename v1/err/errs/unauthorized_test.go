package errs

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorUnauthorized(t *testing.T) {
	tests := []struct {
		name    string
		source  []byte
		id      fields.TracebackID
		want    []byte
		wantErr bool
	}{
		{"nil source", nil, tracebackIDOne, encoded(16, tracebackIDOne), false},
		{"empty source", []byte{}, tracebackIDOne, encoded(16, tracebackIDOne), false},
		{"binary source", []byte{0x01, 0x02}, tracebackIDOne, encoded(16, tracebackIDOne), false},
		{"address source", []byte("192.168.1.1"), tracebackIDMax, encoded(16, tracebackIDMax), false},
		{"zero value", nil, fields.TracebackID{}, encoded(16, fields.TracebackID{}), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorUnauthorizedWithTracebackID(tt.source, tt.id)
			assertProtocolError(t, e, fields.Unauthorized, tt.id, tt.want, tt.wantErr)
		})
	}
}
