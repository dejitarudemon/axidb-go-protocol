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
	got         uint32
	limit       uint32
	tracebackID fields.TracebackID
}

func NewErrorBatchLimitIsExceeded(got, limit uint32) ErrorBatchLimitIsExceeded {
	return ErrorBatchLimitIsExceeded{
		got:         got,
		limit:       limit,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorBatchLimitIsExceededWithTracebackID(got, limit uint32, tracebackID fields.TracebackID) ErrorBatchLimitIsExceeded {
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
	return err.FieldSize + err.TracebackIDFieldSize + CurrentBatchLimitFieldSize
}

func (e ErrorBatchLimitIsExceeded) Code() err.Code {
	return err.BatchLimitIsExceeded
}

func (e ErrorBatchLimitIsExceeded) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
	buf.AppendUint32(e.limit)
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
