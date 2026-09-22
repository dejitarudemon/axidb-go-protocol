package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorInternalError{}

type ErrorInternalError struct {
	err         error
	tracebackID fields.TracebackID
}

func NewErrorInternalError(err error) ErrorInternalError {
	return ErrorInternalError{
		err:         err,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorInternalErrorWithTracebackID(err error, tracebackID fields.TracebackID) ErrorInternalError {
	return ErrorInternalError{
		err:         err,
		tracebackID: tracebackID,
	}
}

func (e ErrorInternalError) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorInternalError) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

func (e ErrorInternalError) Code() fields.Error {
	return fields.InternalError
}

func (e ErrorInternalError) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

func (e ErrorInternalError) Error() string {
	return fmt.Sprintf("%v %v: %v,", e.tracebackID, e.Code(), e.err)
}

func (e ErrorInternalError) Unwrap() error {
	return e.err
}

func (e ErrorInternalError) IsValid() error {
	if e.err == nil {
		return err.NewValidationError(
			"source error is nil",
			"error", "ErrorInternalError",
			"tracebackID", e.tracebackID,
		)
	}
	if er, ok := e.err.(err.ProtocolError); ok {
		return err.NewValidationError(
			"source error is ProtocolError",
			"error", "ErrorInternalError",
			"source", er,
			"tracebackID", e.tracebackID,
		)
	}

	return nil
}
