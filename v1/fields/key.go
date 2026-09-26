package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// Key is an opaque byte sequence identifying a stored record.
type Key []byte

// KeyLenFieldSize is the encoded size in bytes of a length prefix when a key length is written separately.
const KeyLenFieldSize = 4

// Encode writes the raw key bytes into buf without a length prefix.
func (k Key) Encode(buf buffer.Appender) {
	buf.Append(k[:])
}

// Size returns the encoded size in bytes (the key length).
func (k Key) Size() int {
	return len(k)
}
