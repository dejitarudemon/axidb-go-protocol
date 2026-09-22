package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = String("")

const StringLenFieldSize = 4

/*
type String представляет собой строку
в формате UTF-8 строки из спецификации протокола v1.
*/
type String string

func (s String) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(s)))
	buf.AppendString(string(s))
}

func (s String) Size() int {
	return StringLenFieldSize + len(s)
}

func (s String) Type() fields.Type {
	return fields.String
}

func (s String) IsValid() error {
	return nil
}
