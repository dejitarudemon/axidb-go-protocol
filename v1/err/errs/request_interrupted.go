package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorRequestInterrupted{}

type ErrorRequestInterrupted struct {
	interrupted fields.RequestID
	tracebackID uuid.UUID
}

func NewErrorRequestInterrupted(interrupted fields.RequestID) ErrorRequestInterrupted {
	return ErrorRequestInterrupted{
		interrupted: interrupted,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorRequestInterruptedWithTracebackID(interrupted fields.RequestID, tracebackID uuid.UUID) ErrorRequestInterrupted {
	return ErrorRequestInterrupted{
		interrupted: interrupted,
		tracebackID: tracebackID,
	}
}

func (e ErrorRequestInterrupted) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorRequestInterrupted) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorRequestInterrupted) Code() err.Code {
	return err.RequestInterrupted
}

func (e ErrorRequestInterrupted) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	buf.Append(e.tracebackID[:])
}

func (e ErrorRequestInterrupted) Error() string {
	return fmt.Sprintf("%v %v: id %v,", e.tracebackID, e.Code(), e.interrupted)
}

func (e ErrorRequestInterrupted) IsValid() error {
	return nil
}
