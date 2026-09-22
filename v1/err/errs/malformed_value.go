package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorMalformedValue{}

type ErrorMalformedValue struct {
	msg         string
	tracebackID fields.TracebackID
}

func NewErrorMalformedValue(msg string) ErrorMalformedValue {
	return ErrorMalformedValue{
		msg:         msg,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorMalformedValueWithTracebackID(msg string, tracebackID fields.TracebackID) ErrorMalformedValue {
	return ErrorMalformedValue{
		msg:         msg,
		tracebackID: tracebackID,
	}
}

func (e ErrorMalformedValue) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorMalformedValue) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize + len(e.msg)
}

func (e ErrorMalformedValue) Code() fields.Error {
	return fields.MalformedValue
}

func (e ErrorMalformedValue) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
	buf.Append([]byte(e.msg))
}

func (e ErrorMalformedValue) Error() string {
	return fmt.Sprintf("%v %v: %v", e.tracebackID, e.Code(), e.msg)
}

func (e ErrorMalformedValue) IsValid() error {
	return nil
}
