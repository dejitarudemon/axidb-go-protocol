package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.ProtocolError = ErrorBodyLimitIsExceeded{}

const CurrentBodyLimitFieldSize = 4

type ErrorBodyLimitIsExceeded struct {
	got         uint32
	limit       uint32
	tracebackID uuid.UUID
}

func NewErrorBodyLimitIsExceeded(got, limit uint32) ErrorBodyLimitIsExceeded {
	return ErrorBodyLimitIsExceeded{
		got:         got,
		limit:       limit,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorBodyLimitIsExceeded) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorBodyLimitIsExceeded) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize + CurrentBodyLimitFieldSize
}

func (e ErrorBodyLimitIsExceeded) Code() err.Code {
	return err.BodyLimitIsExceeded
}

func (e ErrorBodyLimitIsExceeded) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
	buf.AppendUint32(e.limit)
}

func (e ErrorBodyLimitIsExceeded) Error() string {
	return fmt.Sprintf("%v %v: got %v bytes but limit is %v bytes,", e.tracebackID, e.Code(), e.got, e.limit)
}

func (e ErrorBodyLimitIsExceeded) IsValid() error {
	if e.got <= e.limit {
		return err.NewValidationError(
			"limit is greater than got",
			"error", "ErrorBodyLimitIsExceeded",
			"got", e.got,
			"limit", e.limit,
			"tracebackID", e.tracebackID,
		)
	}
	return nil
}
