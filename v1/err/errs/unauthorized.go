package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnauthorized{}

// ErrorUnauthorized is a protocol error with code [fields.Unauthorized] reporting failed authorization.
type ErrorUnauthorized struct {
	source      []byte
	tracebackID fields.TracebackID
}

// NewErrorUnauthorized returns an ErrorUnauthorized for source with a newly generated traceback ID.
func NewErrorUnauthorized(source []byte) ErrorUnauthorized {
	return ErrorUnauthorized{
		source:      source,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorUnauthorizedWithTracebackID returns an ErrorUnauthorized for source with the given traceback ID.
func NewErrorUnauthorizedWithTracebackID(source []byte, tracebackID fields.TracebackID) ErrorUnauthorized {
	return ErrorUnauthorized{
		source:      source,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorUnauthorized) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorUnauthorized) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.Unauthorized].
func (e ErrorUnauthorized) Code() fields.Error {
	return fields.Unauthorized
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorUnauthorized) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorUnauthorized) Error() string {
	return fmt.Sprintf("%v %v: % X,", e.tracebackID, e.Code(), e.source)
}

// IsValid reports whether the error payload is consistent with its protocol code.
func (e ErrorUnauthorized) IsValid() error {
	return nil
}
