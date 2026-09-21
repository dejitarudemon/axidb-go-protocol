package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type Version предназначен для хранения Version, его валидации и кодирования в сообщение.
*/
type Version uint32

const VersionFieldSize = 1

func (r Version) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(r))
}

func (r Version) Size() int {
	return VersionFieldSize
}
