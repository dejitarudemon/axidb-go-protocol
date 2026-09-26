package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorBodyLimitIsExceeded{}

// CurrentBodyLimitFieldSize is the encoded size in bytes of the body limit field in this error.
const CurrentBodyLimitFieldSize = 4

// ErrorBodyLimitIsExceeded is a protocol error with code [fields.BodyLimitIsExceeded] reporting an oversized body.
type ErrorBodyLimitIsExceeded struct {
	got         uint32
	limit       fields.BodyLimit
	tracebackID fields.TracebackID
}

// NewErrorBodyLimitIsExceeded returns an ErrorBodyLimitIsExceeded with a newly generated traceback ID.
func NewErrorBodyLimitIsExceeded(got uint32, limit fields.BodyLimit) ErrorBodyLimitIsExceeded {
	return ErrorBodyLimitIsExceeded{
		got:         got,
		limit:       limit,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorBodyLimitIsExceededWithTracebackID returns an ErrorBodyLimitIsExceeded with the given traceback ID.
func NewErrorBodyLimitIsExceededWithTracebackID(got uint32, limit fields.BodyLimit, tracebackID fields.TracebackID) ErrorBodyLimitIsExceeded {
	return ErrorBodyLimitIsExceeded{
		got:         got,
		limit:       limit,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorBodyLimitIsExceeded) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorBodyLimitIsExceeded) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize + CurrentBodyLimitFieldSize
}

// Code returns [fields.BodyLimitIsExceeded].
func (e ErrorBodyLimitIsExceeded) Code() fields.Error {
	return fields.BodyLimitIsExceeded
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorBodyLimitIsExceeded) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
	e.limit.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorBodyLimitIsExceeded) Error() string {
	return fmt.Sprintf("%v %v: got %v bytes but limit is %v bytes,", e.tracebackID, e.Code(), e.got, e.limit)
}

// IsValid reports whether the error payload is consistent with its protocol code.
// Got body size must exceed limit.
func (e ErrorBodyLimitIsExceeded) IsValid() error {
	if fields.BodyLimit(e.got) <= e.limit {
		return err.NewValidationError(
			"got is not greater than limit",
			"error", "ErrorBodyLimitIsExceeded",
			"got", e.got,
			"limit", e.limit,
			"tracebackID", e.tracebackID,
		)
	}
	return nil
}
