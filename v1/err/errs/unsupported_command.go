package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnsupportedCommand{}

type ErrorUnsupportedCommand struct {
	got         fields.Command
	tracebackID uuid.UUID
}

func NewErrorUnsupportedCommand(got fields.Command) ErrorUnsupportedCommand {
	return ErrorUnsupportedCommand{
		got:         got,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorUnsupportedCommandWithTracebackID(got fields.Command, tracebackID uuid.UUID) ErrorUnsupportedCommand {
	return ErrorUnsupportedCommand{
		got:         got,
		tracebackID: tracebackID,
	}
}

func (e ErrorUnsupportedCommand) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorUnsupportedCommand) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorUnsupportedCommand) Code() err.Code {
	return err.UnsupportedCommand
}

func (e ErrorUnsupportedCommand) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	buf.Append(e.tracebackID[:])
}

func (e ErrorUnsupportedCommand) Error() string {
	return fmt.Sprintf("%v %v: %v,", e.tracebackID, e.Code(), e.got)
}

func (e ErrorUnsupportedCommand) IsValid() error {
	// Реализующий протокол обязан реализовывать официальные команды.
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
