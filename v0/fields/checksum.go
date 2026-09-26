package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// Checksum is a CRC-32/XFER frame checksum.
//
// Parameters: poly=0x000000AF, init=0, refin=false, refout=false, xorout=0.
type Checksum uint32

// ChecksumFieldSize is the encoded size in bytes of a [Checksum].
const ChecksumFieldSize = 4

const xferPoly uint32 = 0x000000AF

var xferTable = func() (t [256]uint32) {
	for i := range t {
		crc := uint32(i) << 24
		for range 8 {
			if crc&0x80000000 != 0 {
				crc = crc<<1 ^ xferPoly
			} else {
				crc <<= 1
			}
		}
		t[i] = crc
	}
	return t
}()

func update(crc uint32, data []byte) uint32 {
	for _, b := range data {
		crc = xferTable[byte(crc>>24)^b] ^ crc<<8
	}
	return crc
}

// NewChecksum returns the CRC-32/XFER checksum of data.
func NewChecksum(data []byte) Checksum {
	return Checksum(update(0, data))
}

// NewChecksumWithParts returns the checksum of base followed by additional byte parts.
func NewChecksumWithParts(base []byte, parts ...[]byte) Checksum {
	crc := update(0, base)
	for _, p := range parts {
		crc = update(crc, p)
	}
	return Checksum(crc)
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
