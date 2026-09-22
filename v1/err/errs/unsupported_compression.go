package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnsupportedCompression{}

type ErrorUnsupportedCompression struct {
	compression fields.Compression
	tracebackID fields.TracebackID
}

func NewErrorUnsupportedCompression(compression fields.Compression) ErrorUnsupportedCompression {
	return ErrorUnsupportedCompression{
		compression: compression,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorUnsupportedCompressionWithTracebackID(compression fields.Compression, tracebackID fields.TracebackID) ErrorUnsupportedCompression {
	return ErrorUnsupportedCompression{
		compression: compression,
		tracebackID: tracebackID,
	}
}

func (e ErrorUnsupportedCompression) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorUnsupportedCompression) Size() int {
	return e.Code().Size() + fields.TracebackIDFieldSize
}

func (e ErrorUnsupportedCompression) Code() fields.Error {
	return fields.UnsupportedCompression
}

func (e ErrorUnsupportedCompression) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

func (e ErrorUnsupportedCompression) Error() string {
	return fmt.Sprintf("%v %v: %v", e.tracebackID, e.Code(), e.compression)
}

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
