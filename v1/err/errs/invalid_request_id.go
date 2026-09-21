package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.ProtocolError = ErrorInvalidRequestID{}

type ErrorInvalidRequestID struct {
	command     command.Code
	requestID   uint32
	tracebackID uuid.UUID
}

func NewErrorInvalidRequestID(command command.Code, requestID uint32) ErrorInvalidRequestID {
	return ErrorInvalidRequestID{
		command:     command,
		requestID:   requestID,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorInvalidRequestID) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorInvalidRequestID) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorInvalidRequestID) Code() err.Code {
	return err.InvalidRequestID
}

func (e ErrorInvalidRequestID) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
}

func (e ErrorInvalidRequestID) Error() string {
	return fmt.Sprintf("%v %v: invalid %v command for request id %v,", e.tracebackID, e.Code(), e.command, e.requestID)
}

// не проверяем команду, т.к. может быть кастомная
func (e ErrorInvalidRequestID) IsValid() error {
	if e.command == command.Handshake && e.requestID == 0 || e.command != command.Answer && e.requestID != 0 {
		return err.NewValidationError(
			"the request id is valid for the command",
			"error", "ErrorInvalidRequestID",
			"command", e.command,
			"requestID", e.requestID,
			"tracebackID", e.tracebackID,
		)
	}

	return nil
}
