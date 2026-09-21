package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type BatchLimit предназначен для хранения BatchLimit, его валидации и кодирования в сообщение.
*/
type BatchLimit uint32

const BatchLimitFieldSize = 4

func (b BatchLimit) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(b))
}

func (b BatchLimit) Size() int {
	return BatchLimitFieldSize
}
