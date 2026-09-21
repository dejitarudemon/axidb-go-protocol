package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.ProtocolError = ErrorUnsupportedCommand{}

type ErrorUnsupportedCommand struct {
	got         command.Code
	tracebackID uuid.UUID
}

func NewErrorUnsupportedCommand(got command.Code) ErrorUnsupportedCommand {
	return ErrorUnsupportedCommand{
		got:         got,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorUnsupportedCommandWithTracebackID(got command.Code, tracebackID uuid.UUID) ErrorUnsupportedCommand {
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
	buf.AppendUint16(uint16(e.Code()))
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
