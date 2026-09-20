package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.Error = ErrorMismatchedChecksum{}

type ErrorMismatchedChecksum struct {
	got         uint32
	expected    uint32
	tracebackID uuid.UUID
}

func NewErrorMismatchedChecksum(got, expected uint32) ErrorMismatchedChecksum {
	return ErrorMismatchedChecksum{
		got:         got,
		expected:    expected,
		tracebackID: generateNewTracebackID(),
	}
}

func (e ErrorMismatchedChecksum) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorMismatchedChecksum) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorMismatchedChecksum) Code() err.Code {
	return err.MismatchedChecksum
}

func (e ErrorMismatchedChecksum) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
}

func (e ErrorMismatchedChecksum) Error() string {
	return fmt.Sprintf("%v %v: got %08X, expected %08X,", e.tracebackID, e.Code(), e.got, e.expected)
}

func (e ErrorMismatchedChecksum) IsValid() bool {
	return e.expected != e.got
}
