package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compression"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.Error = ErrorProhibitedCompression{}

type ErrorProhibitedCompression struct {
	compression compression.Code
	command     command.Code
	tracebackID uuid.UUID
}

func NewErrorProhibitedCompression(compression compression.Code, command command.Code) ErrorProhibitedCompression {
	return ErrorProhibitedCompression{
		compression: compression,
		command:     command,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorProhibitedCompression) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorProhibitedCompression) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorProhibitedCompression) Code() err.Code {
	return err.ProhibitedCompression
}

func (e ErrorProhibitedCompression) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
}

func (e ErrorProhibitedCompression) Error() string {
	return fmt.Sprintf("%v %v: %v compression for %v command", e.tracebackID, e.Code(), e.compression, e.command)
}
