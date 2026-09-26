package errs

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorBodyLimitIsExceeded(t *testing.T) {
	tests := []struct {
		name    string
		got     uint32
		limit   fields.BodyLimit
		id      fields.TracebackID
		want    []byte
		wantErr bool
	}{
		{"got below limit", 0, 1, tracebackIDOne, encoded(6, tracebackIDOne, 0x00, 0x00, 0x00, 0x01), true},
		{"got above limit", 1, 0, tracebackIDOne, encoded(6, tracebackIDOne, 0x00, 0x00, 0x00, 0x00), false},
		{"got equals limit", 1, 1, tracebackIDOne, encoded(6, tracebackIDOne, 0x00, 0x00, 0x00, 0x01), true},
		{"max limit", 1, math.MaxUint32, tracebackIDMax, encoded(6, tracebackIDMax, 0xFF, 0xFF, 0xFF, 0xFF), true},
		{"max got", math.MaxUint32, 1, tracebackIDOne, encoded(6, tracebackIDOne, 0x00, 0x00, 0x00, 0x01), false},
		{"max got equals max limit", math.MaxUint32, math.MaxUint32, tracebackIDOne, encoded(6, tracebackIDOne, 0xFF, 0xFF, 0xFF, 0xFF), true},
		{"zero value", 0, 0, fields.TracebackID{}, encoded(6, fields.TracebackID{}, 0x00, 0x00, 0x00, 0x00), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorBodyLimitIsExceededWithTracebackID(tt.got, tt.limit, tt.id)
			assertProtocolError(t, e, fields.BodyLimitIsExceeded, tt.id, tt.want, tt.wantErr)
		})
	}
}
