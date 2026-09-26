package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorMalformedValue{}

// ErrorMalformedValue is a protocol error with code [fields.MalformedValue] reporting an invalid field value.
type ErrorMalformedValue struct {
	msg         string
	tracebackID fields.TracebackID
}

// NewErrorMalformedValue returns an ErrorMalformedValue for msg with a newly generated traceback ID.
func NewErrorMalformedValue(msg string) ErrorMalformedValue {
	return ErrorMalformedValue{
		msg:         msg,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorMalformedValueWithTracebackID returns an ErrorMalformedValue for msg with the given traceback ID.
func NewErrorMalformedValueWithTracebackID(msg string, tracebackID fields.TracebackID) ErrorMalformedValue {
	return ErrorMalformedValue{
		msg:         msg,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorMalformedValue) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorMalformedValue) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize + len(e.msg)
}

// Code returns [fields.MalformedValue].
func (e ErrorMalformedValue) Code() fields.Error {
	return fields.MalformedValue
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorMalformedValue) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
	buf.Append([]byte(e.msg))
}

// Error returns a human-readable summary of the error.
func (e ErrorMalformedValue) Error() string {
	return fmt.Sprintf("%v %v: %v", e.tracebackID, e.Code(), e.msg)
}

// IsValid reports whether the error payload is consistent with its protocol code.
func (e ErrorMalformedValue) IsValid() error {
	return nil
}
