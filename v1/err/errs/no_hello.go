package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorNoHello{}

// ErrorNoHello is a protocol error with code [fields.NoHello] reporting that a hello handshake was not received.
type ErrorNoHello struct {
	tracebackID fields.TracebackID
}

// NewErrorNoHello returns an ErrorNoHello with a newly generated traceback ID.
func NewErrorNoHello() ErrorNoHello {
	return ErrorNoHello{
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorNoHelloWithTracebackID returns an ErrorNoHello with the given traceback ID.
func NewErrorNoHelloWithTracebackID(tracebackID fields.TracebackID) ErrorNoHello {
	return ErrorNoHello{
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorNoHello) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorNoHello) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.NoHello].
func (e ErrorNoHello) Code() fields.Error {
	return fields.NoHello
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorNoHello) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorNoHello) Error() string {
	return fmt.Sprintf("%v %v,", e.tracebackID, e.Code())
}

// IsValid reports whether the error payload is consistent with its protocol code.
func (e ErrorNoHello) IsValid() error {
	return nil
}
