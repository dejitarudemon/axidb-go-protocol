package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = Int(0)

// IntValueFieldSize is the encoded size in bytes of an [Int] payload.
const IntValueFieldSize = 8

// Int is a signed 64-bit integer value with type code [fields.Int].
type Int int64

// Encode writes the wire encoding of the value into buf.
func (i Int) Encode(buf buffer.Appender) {
	buf.AppendUint64(uint64(i))
}

// Size returns the encoded value size in bytes.
func (i Int) Size() int {
	return IntValueFieldSize
}

// Type returns [fields.Int].
func (i Int) Type() fields.Type {
	return fields.Int
}

// IsValid reports whether the value satisfies type-specific rules.
func (i Int) IsValid() error {
	return nil
}
