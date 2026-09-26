package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnsupportedCommand{}

// ErrorUnsupportedCommand is a protocol error with code [fields.UnsupportedCommand] reporting an unsupported command.
type ErrorUnsupportedCommand struct {
	got         fields.Command
	tracebackID fields.TracebackID
}

// NewErrorUnsupportedCommand returns an ErrorUnsupportedCommand for got with a newly generated traceback ID.
func NewErrorUnsupportedCommand(got fields.Command) ErrorUnsupportedCommand {
	return ErrorUnsupportedCommand{
		got:         got,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorUnsupportedCommandWithTracebackID returns an ErrorUnsupportedCommand for got with the given traceback ID.
func NewErrorUnsupportedCommandWithTracebackID(got fields.Command, tracebackID fields.TracebackID) ErrorUnsupportedCommand {
	return ErrorUnsupportedCommand{
		got:         got,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorUnsupportedCommand) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorUnsupportedCommand) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.UnsupportedCommand].
func (e ErrorUnsupportedCommand) Code() fields.Error {
	return fields.UnsupportedCommand
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorUnsupportedCommand) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorUnsupportedCommand) Error() string {
	return fmt.Sprintf("%v %v: %v,", e.tracebackID, e.Code(), e.got)
}

// IsValid reports whether the error payload is consistent with its protocol code.
// Official protocol commands must be supported by the implementation.
func (e ErrorUnsupportedCommand) IsValid() error {
	if e.got.IsValid() {
		return err.NewValidationError(
			"command must be supported",
			"error", "ErrorUnsupportedCommand",
			"got", e.got,
			"tracebackID", e.tracebackID,
		)
	}
	return nil
}
