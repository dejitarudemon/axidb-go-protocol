package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorInternalError{}

// ErrorInternalError is a protocol error with code [fields.InternalError] wrapping a non-protocol failure.
type ErrorInternalError struct {
	err         error
	tracebackID fields.TracebackID
}

// NewErrorInternalError returns an ErrorInternalError for err with a newly generated traceback ID.
func NewErrorInternalError(err error) ErrorInternalError {
	return ErrorInternalError{
		err:         err,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorInternalErrorWithTracebackID returns an ErrorInternalError for err with the given traceback ID.
func NewErrorInternalErrorWithTracebackID(err error, tracebackID fields.TracebackID) ErrorInternalError {
	return ErrorInternalError{
		err:         err,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorInternalError) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorInternalError) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.InternalError].
func (e ErrorInternalError) Code() fields.Error {
	return fields.InternalError
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorInternalError) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorInternalError) Error() string {
	return fmt.Sprintf("%v %v: %v,", e.tracebackID, e.Code(), e.err)
}

// Unwrap returns the wrapped source error.
func (e ErrorInternalError) Unwrap() error {
	return e.err
}

// IsValid reports whether the error payload is consistent with its protocol code.
// The source error must be non-nil and must not itself be a [err.ProtocolError].
func (e ErrorInternalError) IsValid() error {
	if e.err == nil {
		return err.NewValidationError(
			"source error is nil",
			"error", "ErrorInternalError",
			"tracebackID", e.tracebackID,
		)
	}
	if er, ok := e.err.(err.ProtocolError); ok {
		return err.NewValidationError(
			"source error is ProtocolError",
			"error", "ErrorInternalError",
			"source", er,
			"tracebackID", e.tracebackID,
		)
	}

	return nil
}
