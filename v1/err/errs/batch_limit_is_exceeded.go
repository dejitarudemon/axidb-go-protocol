package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorBatchLimitIsExceeded{}

// CurrentBatchLimitFieldSize is the encoded size in bytes of the batch limit field in this error.
const CurrentBatchLimitFieldSize = 4

// ErrorBatchLimitIsExceeded is a protocol error with code [fields.BatchLimitIsExceeded] reporting too many batch requests.
type ErrorBatchLimitIsExceeded struct {
	got         uint32
	limit       fields.BatchLimit
	tracebackID fields.TracebackID
}

// NewErrorBatchLimitIsExceeded returns an ErrorBatchLimitIsExceeded with a newly generated traceback ID.
func NewErrorBatchLimitIsExceeded(got uint32, limit fields.BatchLimit) ErrorBatchLimitIsExceeded {
	return ErrorBatchLimitIsExceeded{
		got:         got,
		limit:       limit,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorBatchLimitIsExceededWithTracebackID returns an ErrorBatchLimitIsExceeded with the given traceback ID.
func NewErrorBatchLimitIsExceededWithTracebackID(got uint32, limit fields.BatchLimit, tracebackID fields.TracebackID) ErrorBatchLimitIsExceeded {
	return ErrorBatchLimitIsExceeded{
		got:         got,
		limit:       limit,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorBatchLimitIsExceeded) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorBatchLimitIsExceeded) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize + CurrentBatchLimitFieldSize
}

// Code returns [fields.BatchLimitIsExceeded].
func (e ErrorBatchLimitIsExceeded) Code() fields.Error {
	return fields.BatchLimitIsExceeded
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorBatchLimitIsExceeded) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
	e.limit.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorBatchLimitIsExceeded) Error() string {
	return fmt.Sprintf("%v %v: got %v requests but limit is %v requests", e.tracebackID, e.Code(), e.got, e.limit)
}

// IsValid reports whether the error payload is consistent with its protocol code.
// Got batch size must exceed limit.
func (e ErrorBatchLimitIsExceeded) IsValid() error {
	if fields.BatchLimit(e.got) <= e.limit {
		return err.NewValidationError(
			"got is not greater than limit",
			"error", "ErrorBatchLimitIsExceeded",
			"got", e.got,
			"limit", e.limit,
			"tracebackID", e.tracebackID,
		)
	}
	return nil
}
