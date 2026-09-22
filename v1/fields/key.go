package fields

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type Key предназначен для хранения Key, его валидации и кодирования в сообщение.
*/
type Key []byte

const KeyLenFieldSize = 4

func (k Key) Encode(buf buffer.Appender) {
	buf.Append(k[:])
}

func (k Key) Size() int {
	return len(k)
}
