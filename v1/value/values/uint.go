package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = Uint(0)

// UintValueFieldSize is the encoded size in bytes of a [Uint] payload.
const UintValueFieldSize = 8

// Uint is an unsigned 64-bit integer value with type code [fields.Uint].
type Uint uint64

// Encode writes the wire encoding of the value into buf.
func (u Uint) Encode(buf buffer.Appender) {
	buf.AppendUint64(uint64(u))
}

// Size returns the encoded value size in bytes.
func (u Uint) Size() int {
	return UintValueFieldSize
}

// Type returns [fields.Uint].
func (u Uint) Type() fields.Type {
	return fields.Uint
}

// IsValid reports whether the value satisfies type-specific rules.
func (u Uint) IsValid() error {
	return nil
}
