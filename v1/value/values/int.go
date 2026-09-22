package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = Int(0)

const IntValueFieldSize = 8

/*
type Int представляет собой целое число
из спецификации протокола v1.
*/
type Int int64

func (i Int) Encode(buf buffer.Appender) {
	buf.AppendUint64(uint64(i))
}

func (i Int) Size() int {
	return IntValueFieldSize
}

func (i Int) Type() fields.Type {
	return fields.Int
}

func (i Int) IsValid() error {
	return nil
}
