package fields

import (
	"bursavich.dev/crc/crc32"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type Checksum предназначен для хранения Checksum, его валидации и кодирования в сообщение.
*/
type Checksum uint32

const ChecksumFieldSize = 4

func NewChecksum(data []byte) Checksum {
	return Checksum(crc32.Castagnoli().Checksum(data))
}

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

func (c Checksum) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(c))
}

func (c Checksum) Size() int {
	return ChecksumFieldSize
}

func (c Checksum) Equal(another Checksum) bool {
	return c == another
}
