package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorNotFound{}

// ErrorNotFound is a protocol error with code [fields.NotFound] reporting that a key was not found.
type ErrorNotFound struct {
	key         fields.Key
	tracebackID fields.TracebackID
}

// NewErrorNotFound returns an ErrorNotFound for key with a newly generated traceback ID.
func NewErrorNotFound(key fields.Key) ErrorNotFound {
	return ErrorNotFound{
		key:         key,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorNotFoundWithTracebackID returns an ErrorNotFound for key with the given traceback ID.
func NewErrorNotFoundWithTracebackID(key fields.Key, tracebackID fields.TracebackID) ErrorNotFound {
	return ErrorNotFound{
		key:         key,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorNotFound) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorNotFound) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.NotFound].
func (e ErrorNotFound) Code() fields.Error {
	return fields.NotFound
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorNotFound) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("%v %v: %q,", e.tracebackID, e.Code(), e.key)
}

// IsValid reports whether the error payload is consistent with its protocol code.
// The key must be non-empty.
func (e ErrorNotFound) IsValid() error {
	if len(e.key) == 0 {
		return err.NewValidationError(
			"key is empty",
			"error", "ErrorNotFound",
			"tracebackID", e.tracebackID,
		)
	}
	return nil
}
