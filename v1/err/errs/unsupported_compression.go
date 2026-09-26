package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnsupportedCompression{}

// ErrorUnsupportedCompression is a protocol error with code [fields.UnsupportedCompression] reporting an unsupported compression algorithm.
type ErrorUnsupportedCompression struct {
	compression fields.Compression
	tracebackID fields.TracebackID
}

// NewErrorUnsupportedCompression returns an ErrorUnsupportedCompression with a newly generated traceback ID.
func NewErrorUnsupportedCompression(compression fields.Compression) ErrorUnsupportedCompression {
	return ErrorUnsupportedCompression{
		compression: compression,
		tracebackID: generateNewTracebackID(),
	}
}

// NewErrorUnsupportedCompressionWithTracebackID returns an ErrorUnsupportedCompression with the given traceback ID.
func NewErrorUnsupportedCompressionWithTracebackID(compression fields.Compression, tracebackID fields.TracebackID) ErrorUnsupportedCompression {
	return ErrorUnsupportedCompression{
		compression: compression,
		tracebackID: tracebackID,
	}
}

// TracebackID returns the error traceback ID.
func (e ErrorUnsupportedCompression) TracebackID() fields.TracebackID {
	return e.tracebackID
}

// Size returns the encoded error message size in bytes.
func (e ErrorUnsupportedCompression) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

// Code returns [fields.UnsupportedCompression].
func (e ErrorUnsupportedCompression) Code() fields.Error {
	return fields.UnsupportedCompression
}

// Encode writes the wire encoding of the error into buf.
func (e ErrorUnsupportedCompression) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

// Error returns a human-readable summary of the error.
func (e ErrorUnsupportedCompression) Error() string {
	return fmt.Sprintf("%v %v: %v", e.tracebackID, e.Code(), e.compression)
}

// IsValid reports whether the error payload is consistent with its protocol code.
// Compression must not be [fields.None].
func (e ErrorUnsupportedCompression) IsValid() error {
	if e.compression == fields.None {
		return err.NewValidationError(
			"unsupported no compression",
			"error", "ErrorUnsupportedCompression",
			"tracebackID", e.tracebackID,
		)
	}

	return nil
}
