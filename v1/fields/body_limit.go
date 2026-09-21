package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type BodyLimit предназначен для хранения BodyLimit, его валидации и кодирования в сообщение.
*/
type BodyLimit uint32

const BodyLimitFieldSize = 4

func (b BodyLimit) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(b))
}

func (b BodyLimit) Size() int {
	return BodyLimitFieldSize
}
