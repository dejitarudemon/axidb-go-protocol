package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.Error = ErrorInternalError{}

type ErrorInternalError struct {
	err         error
	tracebackID uuid.UUID
}

func NewErrorInternalError(err error) ErrorInternalError {
	return ErrorInternalError{
		err:         err,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorInternalError) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorInternalError) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorInternalError) Code() err.Code {
	return err.InternalError
}

func (e ErrorInternalError) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
}

func (e ErrorInternalError) Error() string {
	return fmt.Sprintf("%v %v: %v,", e.tracebackID, e.Code(), e.err)
}

func (e ErrorInternalError) Unwrap() error {
	return e.err
}
