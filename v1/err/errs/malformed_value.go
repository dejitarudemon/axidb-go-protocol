package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/google/uuid"
)

var _ err.ProtocolError = ErrorMalformedValue{}

type ErrorMalformedValue struct {
	msg         string
	tracebackID uuid.UUID
}

func NewErrorMalformedValue(msg string) ErrorMalformedValue {
	return ErrorMalformedValue{
		msg:         msg,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorMalformedValueWithTracebackID(msg string, tracebackID uuid.UUID) ErrorMalformedValue {
	return ErrorMalformedValue{
		msg:         msg,
		tracebackID: tracebackID,
	}
}

func (e ErrorMalformedValue) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorMalformedValue) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize + len(e.msg)
}

func (e ErrorMalformedValue) Code() err.Code {
	return err.MalformedValue
}

func (e ErrorMalformedValue) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	buf.Append(e.tracebackID[:])
	buf.Append([]byte(e.msg))
}

func (e ErrorMalformedValue) Error() string {
	return fmt.Sprintf("%v %v: %v", e.tracebackID, e.Code(), e.msg)
}

func (e ErrorMalformedValue) IsValid() error {
	return nil
}
