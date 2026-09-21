package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.ProtocolError = ErrorUnexpectedCommandInBatch{}

const RequestNumberFieldSize = 4

type ErrorUnexpectedCommandInBatch struct {
	command       command.Code
	requestNumber uint32
	tracebackID   uuid.UUID
}

func NewErrorUnexpectedCommandInBatch(command command.Code, requestNumber uint32) ErrorUnexpectedCommandInBatch {
	return ErrorUnexpectedCommandInBatch{
		command:       command,
		requestNumber: requestNumber,
		tracebackID:   generateNewTracebackID(),
	}
}

func NewErrorUnexpectedCommandInBatchWithTracebackID(command command.Code, requestNumber uint32, tracebackID uuid.UUID) ErrorUnexpectedCommandInBatch {
	return ErrorUnexpectedCommandInBatch{
		command:       command,
		requestNumber: requestNumber,
		tracebackID:   tracebackID,
	}
}

func (e ErrorUnexpectedCommandInBatch) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorUnexpectedCommandInBatch) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize + RequestNumberFieldSize
}

func (e ErrorUnexpectedCommandInBatch) Code() err.Code {
	return err.UnexpectedCommandInBatch
}

func (e ErrorUnexpectedCommandInBatch) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
	buf.AppendUint32(e.requestNumber)
}

func (e ErrorUnexpectedCommandInBatch) Error() string {
	return fmt.Sprintf("%v %v: got %v command for request id %v", e.tracebackID, e.Code(), e.command, e.requestNumber)
}

// не проверяем команду на валидность, т.к. может быть кастомная
func (e ErrorUnexpectedCommandInBatch) IsValid() error {
	switch e.command {
	case command.Read, command.Write, command.Delete:
		return err.NewValidationError(
			"command is allowed for a batch",
			"error", "ErrorUnexpectedCommandInBatch",
			"command", e.command,
			"tracebackID", e.tracebackID,
		)
	}

	return nil
}
