package errs

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorMismatchedChecksum(t *testing.T) {
	tests := []struct {
		name     string
		got      fields.Checksum
		expected fields.Checksum
		id       fields.TracebackID
		want     []byte
		wantErr  bool
	}{
		{"got below expected", 0, 1, tracebackIDOne, encoded(7, tracebackIDOne), false},
		{"got above expected", 1, 0, tracebackIDOne, encoded(7, tracebackIDOne), false},
		{"equal checksums", 1, 1, tracebackIDOne, encoded(7, tracebackIDOne), true},
		{"max expected", 1, math.MaxUint32, tracebackIDMax, encoded(7, tracebackIDMax), false},
		{"max got", math.MaxUint32, 1, tracebackIDOne, encoded(7, tracebackIDOne), false},
		{"equal max checksums", math.MaxUint32, math.MaxUint32, tracebackIDOne, encoded(7, tracebackIDOne), true},
		{"zero value", 0, 0, fields.TracebackID{}, encoded(7, fields.TracebackID{}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorMismatchedChecksumWithTracebackID(tt.got, tt.expected, tt.id)
			assertProtocolError(t, e, fields.MismatchedChecksum, tt.id, tt.want, tt.wantErr)
		})
	}
}
