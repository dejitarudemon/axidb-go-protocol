package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorRequestInterrupted{}

// ErrorRequestInterrupted is a protocol error with code [fields.RequestInterrupted] reporting an interrupted request.
type ErrorRequestInterrupted struct {
	interrupted fields.RequestID
	tracebackID fields.TracebackID
}

// NewErrorRequestInterrupted returns an ErrorRequestInterrupted for interrupted with a newly generated traceback ID.
func NewErrorRequestInterrupted(interrupted fields.RequestID) ErrorRequestInterrupted {
	return ErrorRequestInterrupted{
		interrupted: interrupted,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorRequestInterruptedWithTracebackID returns an ErrorRequestInterrupted for interrupted with the given traceback ID.
func NewErrorRequestInterruptedWithTracebackID(interrupted fields.RequestID, tracebackID fields.TracebackID) ErrorRequestInterrupted {
	return ErrorRequestInterrupted{
		interrupted: interrupted,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorRequestInterrupted) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorRequestInterrupted) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.RequestInterrupted].
func (e ErrorRequestInterrupted) Code() fields.Error {
	return fields.RequestInterrupted
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorRequestInterrupted) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorRequestInterrupted) Error() string {
	return fmt.Sprintf("%v %v: %v,", e.tracebackID, e.Code(), e.interrupted)
}

// IsValid reports whether the error payload is consistent with its protocol code.
func (e ErrorRequestInterrupted) IsValid() error {
	return nil
}
