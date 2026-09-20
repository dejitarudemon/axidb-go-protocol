package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.Error = ErrorInvalidRequestID{}

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

func (e ErrorInvalidRequestID) IsValid() bool {
	if !e.command.IsValid() {
		return false
	}

	switch e.command {
	case command.Answer:
		return true
	case command.Handshake:
		return e.requestID == 0
	}

	return e.requestID != 0
}
