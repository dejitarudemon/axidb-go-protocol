package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type TracebackID предназначен для хранения TracebackID, его валидации и кодирования в сообщение.
*/
type TracebackID [16]byte

const TracebackIDFieldSize = 16

func (t TracebackID) Encode(buf buffer.Appender) {
	buf.Append(t[:])
}

func (t TracebackID) Size() int {
	return TracebackIDFieldSize
}
