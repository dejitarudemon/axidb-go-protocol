package errs

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorMalformedValue(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		id      fields.TracebackID
		want    []byte
		wantErr bool
	}{
		{"empty message", "", tracebackIDOne, encoded(9, tracebackIDOne), false},
		{"ascii message", "a", tracebackIDOne, encoded(9, tracebackIDOne, 'a'), false},
		{"multibyte message", "ф", tracebackIDMax, encoded(9, tracebackIDMax, 0xD1, 0x84), false},
		{"zero value", "", fields.TracebackID{}, encoded(9, fields.TracebackID{}), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorMalformedValueWithTracebackID(tt.msg, tt.id)
			assertProtocolError(t, e, fields.MalformedValue, tt.id, tt.want, tt.wantErr)
		})
	}
}
