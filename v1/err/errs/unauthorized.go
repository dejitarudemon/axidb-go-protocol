package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/google/uuid"
)

var _ err.ProtocolError = ErrorUnauthorized{}

type ErrorUnauthorized struct {
	source      []byte
	tracebackID uuid.UUID
}

func NewErrorUnauthorized(source []byte) ErrorUnauthorized {
	return ErrorUnauthorized{
		source:      source,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorUnauthorized) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorUnauthorized) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorUnauthorized) Code() err.Code {
	return err.Unauthorized
}

func (e ErrorUnauthorized) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
}

func (e ErrorUnauthorized) Error() string {
	return fmt.Sprintf("%v %v: from %v", e.tracebackID, e.Code(), e.source)
}

func (e ErrorUnauthorized) IsValid() error {
	return nil
}
