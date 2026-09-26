package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = Bytes([]byte{})

// BytesLenFieldSize is the encoded size in bytes of the length prefix for [Bytes].
const BytesLenFieldSize = 4

// Bytes is an opaque byte sequence with type code [fields.Bytes].
type Bytes []byte

// Encode writes the wire encoding of the value into buf.
func (b Bytes) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(b)))
	buf.Append(b)
}

// Size returns the encoded value size in bytes.
func (b Bytes) Size() int {
	return BytesLenFieldSize + len(b)
}

// Type returns [fields.Bytes].
func (b Bytes) Type() fields.Type {
	return fields.Bytes
}

// IsValid reports whether the value satisfies type-specific rules.
func (b Bytes) IsValid() error {
	return nil
}
