package errs

import (
	"fmt"
	"uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	err "github.com/dejitarudemon/axidb-go-protocol/v1/error"
)

type ErrorUnsupportedCommand struct {
	got         uint8
	tracebackID uuid.UUID
}

func NewErrorUnsupportedCommand(got uint8) ErrorUnsupportedCommand {
	return ErrorUnsupportedCommand{
		got:         got,
		tracebackID: generateNewTracebackID(),
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
