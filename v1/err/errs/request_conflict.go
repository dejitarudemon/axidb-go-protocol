package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorRequestsConflict{}

// ErrorRequestsConflict is a protocol error with code [fields.RequestsConflict] reporting a duplicate request ID.
type ErrorRequestsConflict struct {
	got         fields.RequestID
	tracebackID fields.TracebackID
}

// NewErrorRequestsConflict returns an ErrorRequestsConflict for got with a newly generated traceback ID.
func NewErrorRequestsConflict(got fields.RequestID) ErrorRequestsConflict {
	return ErrorRequestsConflict{
		got:         got,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorRequestsConflictWithTracebackID returns an ErrorRequestsConflict for got with the given traceback ID.
func NewErrorRequestsConflictWithTracebackID(got fields.RequestID, tracebackID fields.TracebackID) ErrorRequestsConflict {
	return ErrorRequestsConflict{
		got:         got,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorRequestsConflict) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorRequestsConflict) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.RequestsConflict].
func (e ErrorRequestsConflict) Code() fields.Error {
	return fields.RequestsConflict
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorRequestsConflict) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorRequestsConflict) Error() string {
	return fmt.Sprintf("%v %v: %v,", e.tracebackID, e.Code(), e.got)
}

// IsValid reports whether the error payload is consistent with its protocol code.
func (e ErrorRequestsConflict) IsValid() error {
	return nil
}
