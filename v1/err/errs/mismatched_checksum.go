package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorMismatchedChecksum{}

type ErrorMismatchedChecksum struct {
	got         fields.Checksum
	expected    fields.Checksum
	tracebackID fields.TracebackID
}

func NewErrorMismatchedChecksum(got, expected fields.Checksum) ErrorMismatchedChecksum {
	return ErrorMismatchedChecksum{
		got:         got,
		expected:    expected,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorMismatchedChecksumWithTracebackID(got, expected fields.Checksum, tracebackID fields.TracebackID) ErrorMismatchedChecksum {
	return ErrorMismatchedChecksum{
		got:         got,
		expected:    expected,
		tracebackID: tracebackID,
	}
}

func (e ErrorMismatchedChecksum) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorMismatchedChecksum) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

func (e ErrorMismatchedChecksum) Code() fields.Error {
	return fields.MismatchedChecksum
}

func (e ErrorMismatchedChecksum) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

func (e ErrorMismatchedChecksum) Error() string {
	return fmt.Sprintf("%v %v: got %08X, expected %08X,", e.tracebackID, e.Code(), e.got, e.expected)
}

func (e ErrorMismatchedChecksum) IsValid() error {
	if e.expected == e.got {
		return err.NewValidationError(
			"checksums are equal",
			"error", "ErrorMismatchedChecksum",
			"expected", e.expected,
			"got", e.got,
			"tracebackID", e.tracebackID,
		)
	}
	return nil
}
