package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorRequestInterrupted{}

type ErrorRequestInterrupted struct {
	interrupted fields.RequestID
	tracebackID fields.TracebackID
}

func NewErrorRequestInterrupted(interrupted fields.RequestID) ErrorRequestInterrupted {
	return ErrorRequestInterrupted{
		interrupted: interrupted,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorRequestInterruptedWithTracebackID(interrupted fields.RequestID, tracebackID fields.TracebackID) ErrorRequestInterrupted {
	return ErrorRequestInterrupted{
		interrupted: interrupted,
		tracebackID: tracebackID,
	}
}

func (e ErrorRequestInterrupted) TracebackID() fields.TracebackID {
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
	e.tracebackID.Encode(buf)
}

func (e ErrorRequestInterrupted) Error() string {
	return fmt.Sprintf("%v %v: id %v,", e.tracebackID, e.Code(), e.interrupted)
}

func (e ErrorRequestInterrupted) IsValid() error {
	return nil
}
