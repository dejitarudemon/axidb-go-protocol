package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compression"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.Error = ErrorUnsupportedCompression{}

type ErrorUnsupportedCompression struct {
	compression compression.Code
	tracebackID uuid.UUID
}

func NewErrorUnsupportedCompression(compression compression.Code) ErrorUnsupportedCompression {
	return ErrorUnsupportedCompression{
		compression: compression,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorUnsupportedCompression) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorUnsupportedCompression) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorUnsupportedCompression) Code() err.Code {
	return err.UnsupportedCompression
}

func (e ErrorUnsupportedCompression) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
}

func (e ErrorUnsupportedCompression) Error() string {
	return fmt.Sprintf("%v %v: %v", e.tracebackID, e.Code(), e.compression)
}
