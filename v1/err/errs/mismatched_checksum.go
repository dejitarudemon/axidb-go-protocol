package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorMismatchedChecksum{}

// ErrorMismatchedChecksum is a protocol error with code [fields.MismatchedChecksum] reporting unequal checksums.
type ErrorMismatchedChecksum struct {
	got         fields.Checksum
	expected    fields.Checksum
	tracebackID fields.TracebackID
}

// NewErrorMismatchedChecksum returns an ErrorMismatchedChecksum with a newly generated traceback ID.
func NewErrorMismatchedChecksum(got, expected fields.Checksum) ErrorMismatchedChecksum {
	return ErrorMismatchedChecksum{
		got:         got,
		expected:    expected,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorMismatchedChecksumWithTracebackID returns an ErrorMismatchedChecksum with the given traceback ID.
func NewErrorMismatchedChecksumWithTracebackID(got, expected fields.Checksum, tracebackID fields.TracebackID) ErrorMismatchedChecksum {
	return ErrorMismatchedChecksum{
		got:         got,
		expected:    expected,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorMismatchedChecksum) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorMismatchedChecksum) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.MismatchedChecksum].
func (e ErrorMismatchedChecksum) Code() fields.Error {
	return fields.MismatchedChecksum
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorMismatchedChecksum) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorMismatchedChecksum) Error() string {
	return fmt.Sprintf("%v %v: got %08X, expected %08X,", e.tracebackID, e.Code(), e.got, e.expected)
}

// IsValid reports whether the error payload is consistent with its protocol code.
// Got and expected checksums must differ.
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
