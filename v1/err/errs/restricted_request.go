package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorRestrictedRequest{}

// ErrorRestrictedRequest is a protocol error with code [fields.RestrictedRequest] reporting a disallowed request.
type ErrorRestrictedRequest struct {
	id          fields.Key
	source      fields.Key
	key         fields.Key
	command     fields.Command
	requestID   fields.RequestID
	tracebackID fields.TracebackID
}

// NewErrorRestrictedRequest returns an ErrorRestrictedRequest with a newly generated traceback ID.
func NewErrorRestrictedRequest(id, source, key fields.Key, command fields.Command, requestID fields.RequestID) ErrorRestrictedRequest {
	return ErrorRestrictedRequest{
		id:          id,
		source:      source,
		key:         key,
		command:     command,
		requestID:   requestID,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorRestrictedRequestWithTracebackID returns an ErrorRestrictedRequest with the given traceback ID.
func NewErrorRestrictedRequestWithTracebackID(id, source, key fields.Key, command fields.Command, requestID fields.RequestID, tracebackID fields.TracebackID) ErrorRestrictedRequest {
	return ErrorRestrictedRequest{
		id:          id,
		source:      source,
		key:         key,
		command:     command,
		requestID:   requestID,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorRestrictedRequest) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorRestrictedRequest) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.RestrictedRequest].
func (e ErrorRestrictedRequest) Code() fields.Error {
	return fields.RestrictedRequest
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorRestrictedRequest) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorRestrictedRequest) Error() string {
	return fmt.Sprintf("%v %v: someone (id %s) from %s has tried %v (key: %q) via request id %v", e.tracebackID, e.Code(), e.id, e.source, e.command, e.key, e.requestID)
}

// IsValid reports whether the error payload is consistent with its protocol code.
func (e ErrorRestrictedRequest) IsValid() error {
	return nil
}
