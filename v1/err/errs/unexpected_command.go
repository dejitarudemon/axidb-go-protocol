package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnexpectedCommand{}

// ErrorUnexpectedCommand is a protocol error with code [fields.UnexpectedCommand] reporting a command mismatch.
type ErrorUnexpectedCommand struct {
	got         fields.Command
	expected    fields.Command
	tracebackID fields.TracebackID
}

// NewErrorUnexpectedCommand returns an ErrorUnexpectedCommand with a newly generated traceback ID.
func NewErrorUnexpectedCommand(got, expected fields.Command) ErrorUnexpectedCommand {
	return ErrorUnexpectedCommand{
		got:         got,
		expected:    expected,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorUnexpectedCommandWithTracebackID returns an ErrorUnexpectedCommand with the given traceback ID.
func NewErrorUnexpectedCommandWithTracebackID(got, expected fields.Command, tracebackID fields.TracebackID) ErrorUnexpectedCommand {
	return ErrorUnexpectedCommand{
		got:         got,
		expected:    expected,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorUnexpectedCommand) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorUnexpectedCommand) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize + fields.CommandFieldSize
}

// Code returns [fields.UnexpectedCommand].
func (e ErrorUnexpectedCommand) Code() fields.Error {
	return fields.UnexpectedCommand
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorUnexpectedCommand) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
	e.expected.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorUnexpectedCommand) Error() string {
	return fmt.Sprintf("%v %v: got %v, expected %v", e.tracebackID, e.Code(), e.got, e.expected)
}

// IsValid reports whether the error payload is consistent with its protocol code.
// Got and expected commands must differ.
func (e ErrorUnexpectedCommand) IsValid() error {
	if e.expected == e.got {
		return err.NewValidationError(
			"commands are equal",
			"error", "ErrorUnexpectedCommand",
			"expected", e.expected,
			"got", e.got,
			"tracebackID", e.tracebackID,
		)
	}
	return nil
}
