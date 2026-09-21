package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.ProtocolError = ErrorNotFound{}

type ErrorNotFound struct {
	key         []byte
	tracebackID uuid.UUID
}

func NewErrorNotFound(key []byte) ErrorNotFound {
	return ErrorNotFound{
		key:         key,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorNotFound) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorNotFound) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorNotFound) Code() err.Code {
	return err.NotFound
}

func (e ErrorNotFound) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
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
