package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnsupportedVersion{}

// ErrorUnsupportedVersion is a protocol error with code [fields.UnsupportedVersion] reporting an unsupported protocol version.
type ErrorUnsupportedVersion struct {
	got         fields.Version
	tracebackID fields.TracebackID
}

// NewErrorUnsupportedVersion returns an ErrorUnsupportedVersion for got with a newly generated traceback ID.
func NewErrorUnsupportedVersion(got fields.Version) ErrorUnsupportedVersion {
	return ErrorUnsupportedVersion{
		got:         got,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorUnsupportedVersionWithTracebackID returns an ErrorUnsupportedVersion for got with the given traceback ID.
func NewErrorUnsupportedVersionWithTracebackID(got fields.Version, tracebackID fields.TracebackID) ErrorUnsupportedVersion {
	return ErrorUnsupportedVersion{
		got:         got,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorUnsupportedVersion) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorUnsupportedVersion) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.UnsupportedVersion].
func (e ErrorUnsupportedVersion) Code() fields.Error {
	return fields.UnsupportedVersion
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorUnsupportedVersion) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorUnsupportedVersion) Error() string {
	return fmt.Sprintf("%v %v: %v,", e.tracebackID, e.Code(), e.got)
}

// IsValid reports whether the error payload is consistent with its protocol code.
func (e ErrorUnsupportedVersion) IsValid() error {
	return nil
}
