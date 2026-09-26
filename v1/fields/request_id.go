package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// RequestID identifies a client request in a frame header or body.
type RequestID uint32

// RequestIDFieldSize is the encoded size in bytes of a [RequestID].
const RequestIDFieldSize = 4

// Encode writes the wire encoding of the request ID into buf.
func (r RequestID) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(r))
}

// Size returns the encoded size in bytes.
func (r RequestID) Size() int {
	return RequestIDFieldSize
}
