package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnsupportedCommand{}

type ErrorUnsupportedCommand struct {
	got         fields.Command
	tracebackID fields.TracebackID
}

func NewErrorUnsupportedCommand(got fields.Command) ErrorUnsupportedCommand {
	return ErrorUnsupportedCommand{
		got:         got,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorUnsupportedCommandWithTracebackID(got fields.Command, tracebackID fields.TracebackID) ErrorUnsupportedCommand {
	return ErrorUnsupportedCommand{
		got:         got,
		tracebackID: tracebackID,
	}
}

func (e ErrorUnsupportedCommand) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorUnsupportedCommand) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

func (e ErrorUnsupportedCommand) Code() fields.Error {
	return fields.UnsupportedCommand
}

func (e ErrorUnsupportedCommand) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
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
