package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorProhibitedCompression{}

// ErrorProhibitedCompression is a protocol error with code [fields.ProhibitedCompression] reporting disallowed compression for a command.
type ErrorProhibitedCompression struct {
	compression fields.Compression
	command     fields.Command
	tracebackID fields.TracebackID
}

// NewErrorProhibitedCompression returns an ErrorProhibitedCompression with a newly generated traceback ID.
func NewErrorProhibitedCompression(compression fields.Compression, command fields.Command) ErrorProhibitedCompression {
	return ErrorProhibitedCompression{
		compression: compression,
		command:     command,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorProhibitedCompressionWithTracebackID returns an ErrorProhibitedCompression with the given traceback ID.
func NewErrorProhibitedCompressionWithTracebackID(compression fields.Compression, command fields.Command, tracebackID fields.TracebackID) ErrorProhibitedCompression {
	return ErrorProhibitedCompression{
		compression: compression,
		command:     command,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorProhibitedCompression) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorProhibitedCompression) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.ProhibitedCompression].
func (e ErrorProhibitedCompression) Code() fields.Error {
	return fields.ProhibitedCompression
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorProhibitedCompression) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorProhibitedCompression) Error() string {
	return fmt.Sprintf("%v %v: %v compression for %v command", e.tracebackID, e.Code(), e.compression, e.command)
}

// IsValid reports whether the error payload is consistent with its protocol code.
// Command and compression validity are not checked because custom values are allowed.
func (e ErrorProhibitedCompression) IsValid() error {
	if e.compression == fields.None {
		return err.NewValidationError(
			"no compression",
			"error", "ErrorProhibitedCompression",
			"tracebackID", e.tracebackID,
		)
	}

	// Handshake, Ping, and Answer responses must be uncompressed.
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
