package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.Error = ErrorRequestInterrupted{}

type ErrorRequestInterrupted struct {
	interrupted uint32
	tracebackID uuid.UUID
}

func NewErrorRequestInterrupted(interrupted uint32) ErrorRequestInterrupted {
	return ErrorRequestInterrupted{
		interrupted: interrupted,
		tracebackID: generateNewTracebackID(),
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
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
}

func (e ErrorRequestInterrupted) Error() string {
	return fmt.Sprintf("%v %v: id %v,", e.tracebackID, e.Code(), e.interrupted)
}

func (e ErrorRequestInterrupted) IsValid() bool {
	return true
}
