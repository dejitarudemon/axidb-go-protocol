package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

type ErrorUnexpectedCommand struct {
	got         command.Code
	expected    command.Code
	tracebackID uuid.UUID
}

func NewErrorUnexpectedCommand(got, expected command.Code) ErrorUnexpectedCommand {
	return ErrorUnexpectedCommand{
		got:         got,
		expected:    expected,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorUnexpectedCommand) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorUnexpectedCommand) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize + command.FieldSize
}

func (e ErrorUnexpectedCommand) Code() err.Code {
	return err.UnexpectedCommand
}

func (e ErrorUnexpectedCommand) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
	e.expected.Encode(buf)
}

func (e ErrorUnexpectedCommand) Error() string {
	return fmt.Sprintf("%v %v: got %v, expected %v", e.tracebackID, e.Code(), e.got, e.expected)
}
