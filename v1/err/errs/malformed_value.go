package errs

import (
	"fmt"
	"uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	err "github.com/dejitarudemon/axidb-go-protocol/v1/error"
)

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
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
	buf.Append([]byte(e.msg))
}

func (e ErrorMalformedValue) Error() string {
	return fmt.Sprintf("%v %v: %v", e.tracebackID, e.Code(), e.msg)
}
