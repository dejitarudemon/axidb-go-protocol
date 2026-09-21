package errs

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ err.ProtocolError = ErrorUnsupportedVersion{}

type ErrorUnsupportedVersion struct {
	got         fields.Version
	tracebackID fields.TracebackID
}

func NewErrorUnsupportedVersion(got fields.Version) ErrorUnsupportedVersion {
	return ErrorUnsupportedVersion{
		got:         got,
		tracebackID: generateNewTracebackID(),
	}
}

func NewErrorUnsupportedVersionWithTracebackID(got fields.Version, tracebackID fields.TracebackID) ErrorUnsupportedVersion {
	return ErrorUnsupportedVersion{
		got:         got,
		tracebackID: tracebackID,
	}
}

func (e ErrorUnsupportedVersion) TracebackID() fields.TracebackID {
	return e.tracebackID
}

func (e ErrorUnsupportedVersion) Size() int {
	return err.FieldSize + err.TracebackIDFieldSize
}

func (e ErrorUnsupportedVersion) Code() err.Code {
	return err.UnsupportedVersion
}

func (e ErrorUnsupportedVersion) Encode(buf buffer.Appender) {
	e.Code().Encode(buf)
	e.tracebackID.Encode(buf)
}

func (e ErrorUnsupportedVersion) Error() string {
	return fmt.Sprintf("%v %v: %v,", e.tracebackID, e.Code(), e.got)
}

func (e ErrorUnsupportedVersion) IsValid() error {
	return nil
}
