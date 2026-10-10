package errs

import (
	"fmt"

	"github.com/dejitarudemon/ignicula-wire/v1/buffer"
	"github.com/dejitarudemon/ignicula-wire/v1/err"
	"github.com/dejitarudemon/ignicula-wire/v1/fields"
)

var _ err.ProtocolError = ErrorCommandNotImplemented{}

// ErrorCommandNotImplemented is a protocol error with code [fields.CommandNotImplemented]
// reporting that a valid command is not implemented by this node.
type ErrorCommandNotImplemented struct {
	tracebackID fields.TracebackID
}

// NewErrorCommandNotImplemented returns an ErrorCommandNotImplemented with a newly generated traceback ID.
func NewErrorCommandNotImplemented() ErrorCommandNotImplemented {
	return ErrorCommandNotImplemented{
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorCommandNotImplementedWithTracebackID returns an ErrorCommandNotImplemented with the given traceback ID.
func NewErrorCommandNotImplementedWithTracebackID(tracebackID fields.TracebackID) ErrorCommandNotImplemented {
	return ErrorCommandNotImplemented{
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorCommandNotImplemented) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorCommandNotImplemented) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.CommandNotImplemented].
func (e ErrorCommandNotImplemented) Code() fields.Error {
	return fields.CommandNotImplemented
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorCommandNotImplemented) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorCommandNotImplemented) Error() string {
	return fmt.Sprintf("%v %v,", e.tracebackID, e.Code())
}

// IsValid reports whether the error payload is consistent with its protocol code.
func (e ErrorCommandNotImplemented) IsValid() error {
	return nil
}
