package values

import (
	"math"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = Float(0)

// FloatValueFieldSize is the encoded size in bytes of a [Float] payload.
const FloatValueFieldSize = 8

// Float is an IEEE 754 binary64 floating-point value with type code [fields.Float].
type Float float64

// Encode writes the wire encoding of the value into buf.
func (f Float) Encode(buf buffer.Appender) {
	buf.AppendUint64(math.Float64bits(float64(f)))
}

// Size returns the encoded value size in bytes.
func (f Float) Size() int {
	return FloatValueFieldSize
}

// Type returns [fields.Float].
func (f Float) Type() fields.Type {
	return fields.Float
}

// IsValid reports whether the value satisfies type-specific rules.
func (f Float) IsValid() error {
	return nil
}
