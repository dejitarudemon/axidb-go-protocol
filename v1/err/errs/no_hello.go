package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorNoHello{}

type ErrorNoHello struct {
	tracebackID fields.TracebackID
}

func NewErrorNoHello() ErrorNoHello {
	return ErrorNoHello{
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorNoHelloWithTracebackID(tracebackID fields.TracebackID) ErrorNoHello {
	return ErrorNoHello{
		tracebackID: tracebackID,
	}
}

func (e ErrorNoHello) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorNoHello) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

func (e ErrorNoHello) Code() fields.Error {
	return fields.NoHello
}

func (e ErrorNoHello) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

func (e ErrorNoHello) Error() string {
	return fmt.Sprintf("%v %v", e.tracebackID, e.Code())
}

func (e ErrorNoHello) IsValid() error {
	return nil
}
