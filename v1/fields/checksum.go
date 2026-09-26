package fields

import (
	"bursavich.dev/crc/crc32"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// Checksum is a Castagnoli CRC-32 frame checksum.
type Checksum uint32

// ChecksumFieldSize is the encoded size in bytes of a [Checksum].
const ChecksumFieldSize = 4

// NewChecksum returns the checksum of data.
func NewChecksum(data []byte) Checksum {
	return Checksum(crc32.Castagnoli().Checksum(data))
}

// NewChecksumWithParts returns the combined checksum of base and additional byte parts.
func NewChecksumWithParts(base []byte, parts ...[]byte) Checksum {
	checksum := crc32.Castagnoli().Checksum(base)

	for _, p := range parts {
		checksum = crc32.Castagnoli().Combine(
			checksum,
			crc32.Castagnoli().Checksum(p),
			int64(len(p)),
		)
	}

	return Checksum(checksum)
}

// Encode writes the wire encoding of the checksum into buf.
func (c Checksum) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(c))
}

// Size returns the encoded size in bytes.
func (c Checksum) Size() int {
	return ChecksumFieldSize
}

// Equal reports whether c equals another checksum.
func (c Checksum) Equal(another Checksum) bool {
	return c == another
}
