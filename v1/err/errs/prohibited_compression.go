package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compression"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.ProtocolError = ErrorProhibitedCompression{}

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

func NewErrorProhibitedCompressionWithTracebackID(compression compression.Code, command command.Code, tracebackID uuid.UUID) ErrorProhibitedCompression {
	return ErrorProhibitedCompression{
		compression: compression,
		command:     command,
		tracebackID: tracebackID,
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

// не проверяем команду и сжатие, т.к. могут быть кастомные.
func (e ErrorProhibitedCompression) IsValid() error {
	if e.compression == compression.None {
		return err.NewValidationError(
			"no compression",
			"error", "ErrorProhibitedCompression",
			"tracebackID", e.tracebackID,
		)
	}

	// некоторые ответы дб без сжатия.
	switch e.command {
	case command.Handshake, command.Ping, command.Answer:
		return nil
	}

	return err.NewValidationError(
		"compression is allowed for command",
		"error", "ErrorProhibitedCompression",
		"compression", e.compression,
		"command", e.command,
		"tracebackID", e.tracebackID,
	)
}
