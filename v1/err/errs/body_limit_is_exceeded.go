package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorBodyLimitIsExceeded{}

const CurrentBodyLimitFieldSize = 4

type ErrorBodyLimitIsExceeded struct {
	got         uint32
	limit       fields.BodyLimit
	tracebackID fields.TracebackID
}

func NewErrorBodyLimitIsExceeded(got uint32, limit fields.BodyLimit) ErrorBodyLimitIsExceeded {
	return ErrorBodyLimitIsExceeded{
		got:         got,
		limit:       limit,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorBodyLimitIsExceededWithTracebackID(got uint32, limit fields.BodyLimit, tracebackID fields.TracebackID) ErrorBodyLimitIsExceeded {
	return ErrorBodyLimitIsExceeded{
		got:         got,
		limit:       limit,
		tracebackID: tracebackID,
	}
}

func (e ErrorBodyLimitIsExceeded) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorBodyLimitIsExceeded) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize + CurrentBodyLimitFieldSize
}

func (e ErrorBodyLimitIsExceeded) Code() fields.Error {
	return fields.BodyLimitIsExceeded
}

func (e ErrorBodyLimitIsExceeded) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
	e.limit.Encode(buf)
}

func (e ErrorBodyLimitIsExceeded) Error() string {
	return fmt.Sprintf("%v %v: got %v bytes but limit is %v bytes,", e.tracebackID, e.Code(), e.got, e.limit)
}

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
