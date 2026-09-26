package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = String("")

// StringLenFieldSize is the encoded size in bytes of the length prefix for [String].
const StringLenFieldSize = 4

// String is a UTF-8 string value with type code [fields.String].
type String string

// Encode writes the wire encoding of the value into buf.
func (s String) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(s)))
	buf.AppendString(string(s))
}

// Size returns the encoded value size in bytes.
func (s String) Size() int {
	return StringLenFieldSize + len(s)
}

// Type returns [fields.String].
func (s String) Type() fields.Type {
	return fields.String
}

// IsValid reports whether the value satisfies type-specific rules.
func (s String) IsValid() error {
	return nil
}
