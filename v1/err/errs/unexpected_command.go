package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnexpectedCommand{}

type ErrorUnexpectedCommand struct {
	got         fields.Command
	expected    fields.Command
	tracebackID uuid.UUID
}

func NewErrorUnexpectedCommand(got, expected fields.Command) ErrorUnexpectedCommand {
	return ErrorUnexpectedCommand{
		got:         got,
		expected:    expected,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorUnexpectedCommandWithTracebackID(got, expected fields.Command, tracebackID uuid.UUID) ErrorUnexpectedCommand {
	return ErrorUnexpectedCommand{
		got:         got,
		expected:    expected,
		tracebackID: tracebackID,
	}
}

func (e ErrorUnexpectedCommand) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorUnexpectedCommand) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize + fields.CommandFieldSize
}

func (e ErrorUnexpectedCommand) Code() err.Code {
	return err.UnexpectedCommand
}

func (e ErrorUnexpectedCommand) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	buf.Append(e.tracebackID[:])
	e.expected.Encode(buf)
}

func (e ErrorUnexpectedCommand) Error() string {
	return fmt.Sprintf("%v %v: got %v, expected %v", e.tracebackID, e.Code(), e.got, e.expected)
}

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
