package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// VersionLen is the number of version bytes that follow a Hello header.
type VersionLen uint8

// VersionLenFieldSize is the encoded size in bytes of a [VersionLen].
const VersionLenFieldSize = 1

// Encode writes the wire encoding of the version count into buf.
func (v VersionLen) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(v))
}

// Size returns the encoded size in bytes.
func (v VersionLen) Size() int {
	return VersionLenFieldSize
}
