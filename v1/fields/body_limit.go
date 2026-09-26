package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// BodyLimit is a maximum allowed frame body size in bytes.
type BodyLimit uint32

// BodyLimitFieldSize is the encoded size in bytes of a [BodyLimit].
const BodyLimitFieldSize = 4

// Encode writes the wire encoding of the body limit into buf.
func (b BodyLimit) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(b))
}

// Size returns the encoded size in bytes.
func (b BodyLimit) Size() int {
	return BodyLimitFieldSize
}
