package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnexpectedCommandInBatch{}

const RequestNumberFieldSize = 4

type ErrorUnexpectedCommandInBatch struct {
	command       fields.Command
	requestNumber fields.RequestNumber
	tracebackID   fields.TracebackID
}

func NewErrorUnexpectedCommandInBatch(command fields.Command, requestNumber fields.RequestNumber) ErrorUnexpectedCommandInBatch {
	return ErrorUnexpectedCommandInBatch{
		command:       command,
		requestNumber: requestNumber,
		tracebackID:   generateNewTracebackID(),
	}
}

func NewErrorUnexpectedCommandInBatchWithTracebackID(command fields.Command, requestNumber fields.RequestNumber, tracebackID fields.TracebackID) ErrorUnexpectedCommandInBatch {
	return ErrorUnexpectedCommandInBatch{
		command:       command,
		requestNumber: requestNumber,
		tracebackID:   tracebackID,
	}
}

func (e ErrorUnexpectedCommandInBatch) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorUnexpectedCommandInBatch) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize + RequestNumberFieldSize
}

func (e ErrorUnexpectedCommandInBatch) Code() err.Code {
	return err.UnexpectedCommandInBatch
}

func (e ErrorUnexpectedCommandInBatch) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
	e.requestNumber.Encode(buf)
}

func (e ErrorUnexpectedCommandInBatch) Error() string {
	return fmt.Sprintf("%v %v: got %v command for request id %v", e.tracebackID, e.Code(), e.command, e.requestNumber)
}

// не проверяем команду на валидность, т.к. может быть кастомная
func (e ErrorUnexpectedCommandInBatch) IsValid() error {
	switch e.command {
	case fields.Read, fields.Write, fields.Delete:
		return err.NewValidationError(
			"command is allowed for a batch",
			"error", "ErrorUnexpectedCommandInBatch",
			"command", e.command,
			"tracebackID", e.tracebackID,
		)
	}

	return nil
}
