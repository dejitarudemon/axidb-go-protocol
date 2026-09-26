package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Body = Read{}

// Read is a [fields.Read] body holding the key to read.
type Read fields.Key

// Size returns the encoded body size in bytes.
func (r Read) Size() int {
	return fields.Key(r).Size()
}

// Encode writes the wire encoding of the body into buf.
func (r Read) Encode(buf buffer.Appender) {
	fields.Key(r).Encode(buf)
}

// Command returns [fields.Read].
func (r Read) Command() fields.Command {
	return fields.Read
}

// IsValid reports whether the body satisfies protocol rules.
// The key must be non-empty.
func (r Read) IsValid() error {
	if fields.Key(r).Size() == 0 {
		return err.NewValidationError(
			"empty key",
			"target", r.Command(),
		)
	}
	return nil
}
