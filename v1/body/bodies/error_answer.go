package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

// ResultNotOK is the failure flag written at the start of an error answer.
var ResultNotOK = []byte{0x00}

var _ body.Answer = ErrorAnswer{}

// ErrorAnswer is a failed reply carrying a protocol error.
type ErrorAnswer struct {
	// Err is the protocol error to encode.
	Err err.ProtocolError
}

// Size returns the encoded body size in bytes.
func (e ErrorAnswer) Size() int {
	if e.Err == nil {
		return ResultFieldSize
	}
	return ResultFieldSize + e.Err.Size()
}

// Encode writes the failure flag and the protocol error into buf.
func (e ErrorAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultNotOK)
	if e.Err != nil {
		e.Err.Encode(buf)
	}
}

// IsResponseTo returns [fields.Answer].
// Do not call this method. It exists only to implement [body.Answer].
// An ErrorAnswer is unambiguous and is decoded as such.
func (e ErrorAnswer) IsResponseTo() fields.Command {
	return fields.Answer
}

// Command returns [fields.Answer].
func (e ErrorAnswer) Command() fields.Command {
	return fields.Answer
}

// IsValid reports whether the body satisfies protocol rules.
// Err must be non-nil.
func (e ErrorAnswer) IsValid() error {
	if e.Err == nil {
		return err.NewValidationError(
			"expected ProtocolError, but got nil",
			"target", e.Command(),
		)
	}
	return nil
}
