/*
protocol values содержит конретные типы, реализующие интерфейс value.V.
Каждый представленный тип является реализацией типа из спецификации протокола v1.
*/

package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = Bytes([]byte{})

const BytesLenFieldSize = 4

/*
type Bytes представляет собой байтовую последовательность
из спецификации протокола v1.
*/
type Bytes []byte

func (b Bytes) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(b)))
	buf.Append(b)
}

func (b Bytes) Size() int {
	return BytesLenFieldSize + len(b)
}

func (b Bytes) Type() fields.Type {
	return fields.Bytes
}

func (b Bytes) IsValid() error {
	return nil
}
