package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorNotFound{}

type ErrorNotFound struct {
	key         []byte
	tracebackID fields.TracebackID
}

func NewErrorNotFound(key []byte) ErrorNotFound {
	return ErrorNotFound{
		key:         key,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorNotFoundWithTracebackID(key []byte, tracebackID fields.TracebackID) ErrorNotFound {
	return ErrorNotFound{
		key:         key,
		tracebackID: tracebackID,
	}
}

func (e ErrorNotFound) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorNotFound) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorNotFound) Code() err.Code {
	return err.NotFound
}

func (e ErrorNotFound) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("%v %v: %q,", e.tracebackID, e.Code(), e.key)
}

func (e ErrorNotFound) IsValid() error {
	if len(e.key) == 0 {
		return err.NewValidationError(
			"key is empty",
			"error", "ErrorNotFound",
			"tracebackID", e.tracebackID,
		)
	}
	return nil
}
