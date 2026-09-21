package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type RequestID предназначен для хранения RequestID, его валидации и кодирования в сообщение.
*/
type RequestID uint32

const RequestIDFieldSize = 4

func (r RequestID) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(r))
}

func (r RequestID) Size() int {
	return RequestIDFieldSize
}
