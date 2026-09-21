package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnsupportedCompression{}

type ErrorUnsupportedCompression struct {
	compression fields.Compression
	tracebackID uuid.UUID
}

func NewErrorUnsupportedCompression(compression fields.Compression) ErrorUnsupportedCompression {
	return ErrorUnsupportedCompression{
		compression: compression,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorUnsupportedCompressionWithTracebackID(compression fields.Compression, tracebackID uuid.UUID) ErrorUnsupportedCompression {
	return ErrorUnsupportedCompression{
		compression: compression,
		tracebackID: tracebackID,
	}
}

func (e ErrorUnsupportedCompression) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorUnsupportedCompression) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorUnsupportedCompression) Code() err.Code {
	return err.UnsupportedCompression
}

func (e ErrorUnsupportedCompression) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
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
