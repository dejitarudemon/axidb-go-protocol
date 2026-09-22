package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorBatchLimitIsExceeded{}

const CurrentBatchLimitFieldSize = 4

type ErrorBatchLimitIsExceeded struct {
	got         fields.BatchLimit
	limit       fields.BatchLimit
	tracebackID fields.TracebackID
}

func NewErrorBatchLimitIsExceeded(got, limit fields.BatchLimit) ErrorBatchLimitIsExceeded {
	return ErrorBatchLimitIsExceeded{
		got:         got,
		limit:       limit,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorBatchLimitIsExceededWithTracebackID(got, limit fields.BatchLimit, tracebackID fields.TracebackID) ErrorBatchLimitIsExceeded {
	return ErrorBatchLimitIsExceeded{
		got:         got,
		limit:       limit,
		tracebackID: tracebackID,
	}
}

func (e ErrorBatchLimitIsExceeded) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorBatchLimitIsExceeded) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize + CurrentBatchLimitFieldSize
}

func (e ErrorBatchLimitIsExceeded) Code() fields.Error {
	return fields.BatchLimitIsExceeded
}

func (e ErrorBatchLimitIsExceeded) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
	e.limit.Encode(buf)
}

func (e ErrorBatchLimitIsExceeded) Error() string {
	return fmt.Sprintf("%v %v: got %v requests but limit is %v requests", e.tracebackID, e.Code(), e.got, e.limit)
}

func (e ErrorBatchLimitIsExceeded) IsValid() error {
	if e.got <= e.limit {
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
