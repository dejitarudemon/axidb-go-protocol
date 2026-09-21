package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorRequestsConflict{}

type ErrorRequestsConflict struct {
	got         fields.RequestID
	tracebackID uuid.UUID
}

func NewErrorRequestsConflict(got fields.RequestID) ErrorRequestsConflict {
	return ErrorRequestsConflict{
		got:         got,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorRequestsConflictWithTracebackID(got fields.RequestID, tracebackID uuid.UUID) ErrorRequestsConflict {
	return ErrorRequestsConflict{
		got:         got,
		tracebackID: tracebackID,
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
	e.Code().Encode(buf)
	buf.Append(e.tracebackID[:])
}

func (e ErrorRequestsConflict) Error() string {
	return fmt.Sprintf("%v %v: id %v,", e.tracebackID, e.Code(), e.got)
}

func (e ErrorRequestsConflict) IsValid() error {
	return nil
}
