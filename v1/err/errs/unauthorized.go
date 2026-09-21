package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnauthorized{}

type ErrorUnauthorized struct {
	source      []byte
	tracebackID fields.TracebackID
}

func NewErrorUnauthorized(source []byte) ErrorUnauthorized {
	return ErrorUnauthorized{
		source:      source,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorUnauthorizedWithTracebackID(source []byte, tracebackID fields.TracebackID) ErrorUnauthorized {
	return ErrorUnauthorized{
		source:      source,
		tracebackID: tracebackID,
	}
}

func (e ErrorUnauthorized) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorUnauthorized) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorUnauthorized) Code() err.Code {
	return err.Unauthorized
}

func (e ErrorUnauthorized) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

func (e ErrorUnauthorized) Error() string {
	return fmt.Sprintf("%v %v: from %v", e.tracebackID, e.Code(), e.source)
}

func (e ErrorUnauthorized) IsValid() error {
	return nil
}
