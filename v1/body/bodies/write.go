package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ body.Body = Write{}

// Write is a [fields.Write] body storing a value under a key.
type Write struct {
	// Key identifies the record.
	Key fields.Key
	// Value is the payload to store.
	Value value.V
}

// Size returns the encoded body size in bytes.
func (w Write) Size() int {
	if w.Value == nil {
		return w.Key.Size() + fields.KeyLenFieldSize
	}

	return w.Key.Size() + fields.KeyLenFieldSize + w.Value.Size() + w.Value.Type().Size()
}

// Encode writes the wire encoding of the body into buf.
func (w Write) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(w.Key.Size()))
	w.Key.Encode(buf)

	if w.Value != nil {
		w.Value.Type().Encode(buf)
		w.Value.Encode(buf)
	}
}

// Command returns [fields.Write].
func (w Write) Command() fields.Command {
	return fields.Write
}

// IsValid reports whether the body satisfies protocol rules.
// Key must be non-empty, and Value must be non-nil and pass its own IsValid check.
func (w Write) IsValid() error {
	if w.Key.Size() == 0 {
		return err.NewValidationError(
			"empty key",
			"target", w.Command(),
		)
	}

	if w.Value == nil {
		return err.NewValidationError(
			"nil value",
			"target", w.Command(),
		)
	}

	return w.Value.IsValid()
}
