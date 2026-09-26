package value

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

// V is a protocol v1 value that can be validated and encoded in a message.
type V interface {
	// Size returns the encoded value size in bytes.
	Size() int

	// IsValid reports whether the value satisfies type-specific rules.
	// Validation details depend on the concrete implementation.
	IsValid() error

	// Encode writes the wire encoding of the value into buf.
	Encode(buf buffer.Appender)

	// Type returns the protocol type code for the value.
	Type() fields.Type
}
