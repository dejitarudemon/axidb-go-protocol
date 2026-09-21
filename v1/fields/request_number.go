package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type RequestNumber предназначен для хранения RequestNumber, его валидации и кодирования в сообщение.
*/
type RequestNumber uint32

const RequestNumberFieldSize = 4

func (r RequestNumber) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(r))
}

func (r RequestNumber) Size() int {
	return RequestNumberFieldSize
}
