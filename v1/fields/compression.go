package fields

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// Compression is a protocol v1 body compression algorithm code.
type Compression uint8

// Compression algorithms defined by protocol v1.
const (
	None Compression = iota
	Zstd
	S2
)

// CompressionFieldSize is the encoded size in bytes of a [Compression] code.
const CompressionFieldSize = 1

// Encode writes the wire encoding of the compression code into buf.
func (c Compression) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(c))
}

// String returns a human-readable compression name.
func (c Compression) String() string {
	switch c {
	case None:
		return "None"
	case Zstd:
		return "Zstd"
	case S2:
		return "S2"
	}

	return fmt.Sprintf("Unknown (%d)", c)
}

// IsValid reports whether c is a known compression code (0 through [S2]).
func (c Compression) IsValid() bool {
	return c <= S2
}

// Size returns the encoded size in bytes.
func (c Compression) Size() int {
	return CompressionFieldSize
}
