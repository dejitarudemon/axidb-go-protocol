package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// TracebackID is a 16-byte identifier attached to protocol errors for correlation.
type TracebackID [16]byte

// TracebackIDFieldSize is the encoded size in bytes of a [TracebackID].
const TracebackIDFieldSize = 16

// Encode writes the wire encoding of the traceback ID into buf.
func (t TracebackID) Encode(buf buffer.Appender) {
	buf.Append(t[:])
}

// Size returns the encoded size in bytes.
func (t TracebackID) Size() int {
	return TracebackIDFieldSize
}
