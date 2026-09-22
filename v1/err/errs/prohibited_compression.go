package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorProhibitedCompression{}

type ErrorProhibitedCompression struct {
	compression fields.Compression
	command     fields.Command
	tracebackID fields.TracebackID
}

func NewErrorProhibitedCompression(compression fields.Compression, command fields.Command) ErrorProhibitedCompression {
	return ErrorProhibitedCompression{
		compression: compression,
		command:     command,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorProhibitedCompressionWithTracebackID(compression fields.Compression, command fields.Command, tracebackID fields.TracebackID) ErrorProhibitedCompression {
	return ErrorProhibitedCompression{
		compression: compression,
		command:     command,
		tracebackID: tracebackID,
	}
}

func (e ErrorProhibitedCompression) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorProhibitedCompression) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

func (e ErrorProhibitedCompression) Code() fields.Error {
	return fields.ProhibitedCompression
}

func (e ErrorProhibitedCompression) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

func (e ErrorProhibitedCompression) Error() string {
	return fmt.Sprintf("%v %v: %v compression for %v command", e.tracebackID, e.Code(), e.compression, e.command)
}

// не проверяем команду и сжатие, т.к. могут быть кастомные.
func (e ErrorProhibitedCompression) IsValid() error {
	if e.compression == fields.None {
		return err.NewValidationError(
			"no compression",
			"error", "ErrorProhibitedCompression",
			"tracebackID", e.tracebackID,
		)
	}

	// некоторые ответы дб без сжатия.
	switch e.command {
	case fields.Handshake, fields.Ping, fields.Answer:
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
