package errs

import (
	"errors"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorInternalError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		id      fields.TracebackID
		want    []byte
		wantErr bool
	}{
		{"nil source", nil, tracebackIDOne, encoded(8, tracebackIDOne), true},
		{"empty source", errors.New(""), tracebackIDOne, encoded(8, tracebackIDOne), false},
		{"source", errors.New("some-error"), tracebackIDMax, encoded(8, tracebackIDMax), false},
		{"protocol error source", NewErrorUnsupportedCommand(fields.Command(10)), tracebackIDOne, encoded(8, tracebackIDOne), true},
		{"invalid protocol error source", NewErrorUnsupportedCommand(fields.Read), tracebackIDOne, encoded(8, tracebackIDOne), true},
		{"zero value", nil, fields.TracebackID{}, encoded(8, fields.TracebackID{}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorInternalErrorWithTracebackID(tt.err, tt.id)
			assertProtocolError(t, e, fields.InternalError, tt.id, tt.want, tt.wantErr)
		})
	}
}

func TestErrorInternalError_Unwrap(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"nil source", nil},
		{"source", errors.New("some-error")},
		{"protocol error source", NewErrorUnsupportedCommand(fields.Command(10))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewErrorInternalErrorWithTracebackID(tt.err, tracebackIDOne).Unwrap(); got != tt.err {
				t.Errorf("Unwrap() = %v, want %v", got, tt.err)
			}
		})
	}
}
