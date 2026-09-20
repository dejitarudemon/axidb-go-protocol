package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.Error = ErrorBatchLimitIsExceeded{}

const CurrentBatchLimitFieldSize = 4

type ErrorBatchLimitIsExceeded struct {
	got         uint32
	limit       uint32
	tracebackID uuid.UUID
}

func NewErrorBatchLimitIsExceeded(got, limit uint32) ErrorBatchLimitIsExceeded {
	return ErrorBatchLimitIsExceeded{
		got:         got,
		limit:       limit,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorBatchLimitIsExceeded) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorBatchLimitIsExceeded) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize + CurrentBatchLimitFieldSize
}

func (e ErrorBatchLimitIsExceeded) Code() err.Code {
	return err.BatchLimitIsExceeded
}

func (e ErrorBatchLimitIsExceeded) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
	buf.AppendUint32(e.limit)
}

func (e ErrorBatchLimitIsExceeded) Error() string {
	return fmt.Sprintf("%v %v: got %v requests but limit is %v bytes,", e.tracebackID, e.Code(), e.got, e.limit)
}
