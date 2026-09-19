package errs

import (
	"fmt"
	"uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	err "github.com/dejitarudemon/axidb-go-protocol/v1/error"
)

type ErrorUnsupportedCompression struct {
	compression uint8
	tracebackID uuid.UUID
}

func NewErrorUnsupportedCompression(compression uint8) ErrorUnsupportedCompression {
	return ErrorUnsupportedCompression{
		compression: compression,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorUnsupportedCompression) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorUnsupportedCompression) Size() int {
	return err.CodeFieldSize + err.TracebackIDFieldSize
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
