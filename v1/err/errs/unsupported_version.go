package errs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
)

var _ err.ProtocolError = ErrorUnsupportedVersion{}

type ErrorUnsupportedVersion struct {
	got         uint8
	tracebackID uuid.UUID
}

func NewErrorUnsupportedVersion(got uint8) ErrorUnsupportedVersion {
	return ErrorUnsupportedVersion{
		got:         got,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorUnsupportedVersionWithTracebackID(got uint8, tracebackID uuid.UUID) ErrorUnsupportedVersion {
	return ErrorUnsupportedVersion{
		got:         got,
		tracebackID: tracebackID,
	}
}

func (e ErrorUnsupportedVersion) TracebackID() uuid.UUID {
	return e.tracebackID
}

func (e ErrorUnsupportedVersion) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorUnsupportedVersion) Code() err.Code {
	return err.UnsupportedVersion
}

func (e ErrorUnsupportedVersion) Encode(buf buffer.Appender) {
	buf.AppendUint16(uint16(e.Code()))
	buf.Append(e.tracebackID[:])
}

func (e ErrorUnsupportedVersion) Error() string {
	return fmt.Sprintf("%v %v: %v,", e.tracebackID, e.Code(), e.got)
}

func (e ErrorUnsupportedVersion) IsValid() error {
	return nil
}
