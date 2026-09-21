package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/google/uuid"
)

var _ err.ProtocolError = ErrorNoHello{}

type ErrorNoHello struct {
	tracebackID uuid.UUID
}

func NewErrorNoHello() ErrorNoHello {
	return ErrorNoHello{
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorNoHello) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorNoHello) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorNoHello) Code() err.Code {
	return err.NoHello
}

func (e ErrorNoHello) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
}

func (e ErrorNoHello) Error() string {
	return fmt.Sprintf("%v %v", e.tracebackID, e.Code())
}

func (e ErrorNoHello) IsValid() error {
	return nil
}
