package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ body.Answer = ReadAnswer{}

// ReadAnswer is a successful reply to [fields.Read] carrying the stored value.
type ReadAnswer struct {
	// Value is the record payload.
	Value value.V
}

// Size returns the encoded body size in bytes.
func (r ReadAnswer) Size() int {
	if r.Value == nil {
		return ResultFieldSize + r.IsResponseTo().Size()
	}

	return ResultFieldSize + r.Value.Type().Size() + r.Value.Size() + r.IsResponseTo().Size()
}

// Encode writes the wire encoding of the body into buf.
func (r ReadAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	r.IsResponseTo().Encode(buf)

	if r.Value != nil {
		r.Value.Type().Encode(buf)
		r.Value.Encode(buf)
	}

}

// IsResponseTo returns [fields.Read].
func (r ReadAnswer) IsResponseTo() fields.Command {
	return fields.Read
}

// Command returns [fields.Answer].
func (r ReadAnswer) Command() fields.Command {
	return fields.Answer
}

// IsValid reports whether the body satisfies protocol rules.
// Value must be non-nil and pass its own IsValid check.
func (r ReadAnswer) IsValid() error {
	if r.Value == nil {
		return err.NewValidationError(
			"nil value",
			"target", r.Command(),
		)
	}

	return r.Value.IsValid()
}
