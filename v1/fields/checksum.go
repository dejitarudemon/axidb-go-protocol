package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type Checksum предназначен для хранения Checksum, его валидации и кодирования в сообщение.
*/
type Checksum uint32

const ChecksumFieldSize = 4

func (c Checksum) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(c))
}

func (c Checksum) Size() int {
	return ChecksumFieldSize
}
