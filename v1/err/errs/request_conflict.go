package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.Error = ErrorRequestsConflict{}

type ErrorRequestsConflict struct {
	got         uint32
	tracebackID uuid.UUID
}

func NewErrorRequestsConflict(got uint32) ErrorRequestsConflict {
	return ErrorRequestsConflict{
		got:         got,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorRequestsConflict) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorRequestsConflict) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorRequestsConflict) Code() err.Code {
	return err.RequestsConflict
}

func (e ErrorRequestsConflict) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
}

func (e ErrorRequestsConflict) Error() string {
	return fmt.Sprintf("%v %v: id %v,", e.tracebackID, e.Code(), e.got)
}

func (e ErrorRequestsConflict) IsValid() bool {
	return true
}
