package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/google/uuid"
)

var _ err.ProtocolError = ErrorRestrictedRequest{}

type ErrorRestrictedRequest struct {
	id          []byte
	source      []byte
	key         []byte
	command     command.Code
	requestID   uint32
	tracebackID uuid.UUID
}

func NewErrorRestrictedRequest(id, source, key []byte, command command.Code, requestID uint32) ErrorRestrictedRequest {
	return ErrorRestrictedRequest{
		id:          id,
		source:      source,
		key:         key,
		command:     command,
		requestID:   requestID,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorRestrictedRequest) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorRestrictedRequest) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorRestrictedRequest) Code() err.Code {
	return err.RestrictedRequest
}

func (e ErrorRestrictedRequest) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
}

func (e ErrorRestrictedRequest) Error() string {
	return fmt.Sprintf("%v %v: someone (id %s) from %s has tried %v (key: %q) via request id %v", e.tracebackID, e.Code(), e.id, e.source, e.command, e.key, e.requestID)
}

func (e ErrorRestrictedRequest) IsValid() error {
	return nil
}
