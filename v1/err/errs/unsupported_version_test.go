package errs

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorUnsupportedVersion(t *testing.T) {
	tests := []struct {
		name    string
		got     fields.Version
		id      fields.TracebackID
		want    []byte
		wantErr bool
	}{
		{"zero version", 0, tracebackIDOne, encoded(1, tracebackIDOne), false},
		{"current version", 1, tracebackIDOne, encoded(1, tracebackIDOne), false},
		{"max one-byte version", 255, tracebackIDMax, encoded(1, tracebackIDMax), false},
		{"zero value", 0, fields.TracebackID{}, encoded(1, fields.TracebackID{}), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorUnsupportedVersionWithTracebackID(tt.got, tt.id)
			assertProtocolError(t, e, fields.UnsupportedVersion, tt.id, tt.want, tt.wantErr)
		})
	}
}
