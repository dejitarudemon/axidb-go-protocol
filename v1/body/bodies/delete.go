package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Body = Delete{}

// Delete is a [fields.Delete] body holding the key to delete.
type Delete fields.Key

// Size returns the encoded body size in bytes.
func (d Delete) Size() int {
	return fields.Key(d).Size()
}

// Encode writes the wire encoding of the body into buf.
func (d Delete) Encode(buf buffer.Appender) {
	fields.Key(d).Encode(buf)
}

// Command returns [fields.Delete].
func (d Delete) Command() fields.Command {
	return fields.Delete
}

// IsValid reports whether the body satisfies protocol rules.
// The key must be non-empty.
func (d Delete) IsValid() error {
	if fields.Key(d).Size() == 0 {
		return err.NewValidationError(
			"empty key",
			"target", d.Command(),
		)
	}

	return nil
}
