package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorRestrictedRequest{}

type ErrorRestrictedRequest struct {
	id          []byte
	source      []byte
	key         []byte
	command     fields.Command
	requestID   fields.RequestID
	tracebackID fields.TracebackID
}

func NewErrorRestrictedRequest(id, source, key []byte, command fields.Command, requestID fields.RequestID) ErrorRestrictedRequest {
	return ErrorRestrictedRequest{
		id:          id,
		source:      source,
		key:         key,
		command:     command,
		requestID:   requestID,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorRestrictedRequestWithTracebackID(id, source, key []byte, command fields.Command, requestID fields.RequestID, tracebackID fields.TracebackID) ErrorRestrictedRequest {
	return ErrorRestrictedRequest{
		id:          id,
		source:      source,
		key:         key,
		command:     command,
		requestID:   requestID,
		tracebackID: tracebackID,
	}
}

func (e ErrorRestrictedRequest) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorRestrictedRequest) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorRestrictedRequest) Code() err.Code {
	return err.RestrictedRequest
}

func (e ErrorRestrictedRequest) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

func (e ErrorRestrictedRequest) Error() string {
	return fmt.Sprintf("%v %v: someone (id %s) from %s has tried %v (key: %q) via request id %v", e.tracebackID, e.Code(), e.id, e.source, e.command, e.key, e.requestID)
}

func (e ErrorRestrictedRequest) IsValid() error {
	return nil
}
