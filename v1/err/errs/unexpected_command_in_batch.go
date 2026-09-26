package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnexpectedCommandInBatch{}

// RequestNumberFieldSize is the encoded size in bytes of the request number field in this error.
const RequestNumberFieldSize = 4

// ErrorUnexpectedCommandInBatch is a protocol error with code [fields.UnexpectedCommandInBatch] reporting a command not allowed in a batch.
type ErrorUnexpectedCommandInBatch struct {
	command       fields.Command
	requestNumber fields.RequestNumber
	tracebackID   fields.TracebackID
}

// NewErrorUnexpectedCommandInBatch returns an ErrorUnexpectedCommandInBatch with a newly generated traceback ID.
func NewErrorUnexpectedCommandInBatch(command fields.Command, requestNumber fields.RequestNumber) ErrorUnexpectedCommandInBatch {
	return ErrorUnexpectedCommandInBatch{
		command:       command,
		requestNumber: requestNumber,
		tracebackID:   generateNewTracebackID(),
	}
}

// NewErrorUnexpectedCommandInBatchWithTracebackID returns an ErrorUnexpectedCommandInBatch with the given traceback ID.
func NewErrorUnexpectedCommandInBatchWithTracebackID(command fields.Command, requestNumber fields.RequestNumber, tracebackID fields.TracebackID) ErrorUnexpectedCommandInBatch {
	return ErrorUnexpectedCommandInBatch{
		command:       command,
		requestNumber: requestNumber,
		tracebackID:   tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorUnexpectedCommandInBatch) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorUnexpectedCommandInBatch) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize + RequestNumberFieldSize
}

// Code returns [fields.UnexpectedCommandInBatch].
func (e ErrorUnexpectedCommandInBatch) Code() fields.Error {
	return fields.UnexpectedCommandInBatch
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorUnexpectedCommandInBatch) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
	e.requestNumber.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorUnexpectedCommandInBatch) Error() string {
	return fmt.Sprintf("%v %v: got %v command for request id %v", e.tracebackID, e.Code(), e.command, e.requestNumber)
}

// IsValid reports whether the error payload is consistent with its protocol code.
// Command validity is not checked because custom commands are allowed.
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
