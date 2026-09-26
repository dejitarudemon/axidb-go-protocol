package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorInvalidRequestID{}

// ErrorInvalidRequestID is a protocol error with code [fields.InvalidRequestID] reporting an invalid request ID for a command.
type ErrorInvalidRequestID struct {
	command     fields.Command
	requestID   fields.RequestID
	tracebackID fields.TracebackID
}

// NewErrorInvalidRequestID returns an ErrorInvalidRequestID with a newly generated traceback ID.
func NewErrorInvalidRequestID(command fields.Command, requestID fields.RequestID) ErrorInvalidRequestID {
	return ErrorInvalidRequestID{
		command:     command,
		requestID:   requestID,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorInvalidRequestIDWithTracebackID returns an ErrorInvalidRequestID with the given traceback ID.
func NewErrorInvalidRequestIDWithTracebackID(command fields.Command, requestID fields.RequestID, tracebackID fields.TracebackID) ErrorInvalidRequestID {
	return ErrorInvalidRequestID{
		command:     command,
		requestID:   requestID,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorInvalidRequestID) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorInvalidRequestID) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.InvalidRequestID].
func (e ErrorInvalidRequestID) Code() fields.Error {
	return fields.InvalidRequestID
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorInvalidRequestID) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorInvalidRequestID) Error() string {
	return fmt.Sprintf("%v %v: invalid %v command for request id %v,", e.tracebackID, e.Code(), e.command, e.requestID)
}

// IsValid reports whether the error payload is consistent with its protocol code.
// Command validity is not checked because custom commands are allowed.
func (e ErrorInvalidRequestID) IsValid() error {
	if e.command == fields.Answer {
		return nil
	}

	if e.command == fields.Handshake {
		if e.requestID == 0 {
			return err.NewValidationError(
				"the request id is valid for the  command",
				"error", "ErrorInvalidRequestID",
				"command", e.command,
				"requestID", e.requestID,
				"tracebackID", e.tracebackID,
			)
		}
		return nil
	}

	if e.requestID != 0 {
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
