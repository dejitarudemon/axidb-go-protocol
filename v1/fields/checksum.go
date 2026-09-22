package fields

import (
	"hash/crc32"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type Checksum предназначен для хранения Checksum, его валидации и кодирования в сообщение.
*/
type Checksum uint32

var table = crc32.MakeTable(crc32.Castagnoli)

const ChecksumFieldSize = 4

func NewChecksum(data []byte) Checksum {
	return Checksum(crc32.Checksum(data, table))
}

func (c Checksum) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(c))
}

func (c Checksum) Size() int {
	return ChecksumFieldSize
}
