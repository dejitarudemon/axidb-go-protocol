package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = Uint(0)

const UintValueFieldSize = 8

type Uint uint64

func (u Uint) Encode(buf buffer.Appender) {
	buf.AppendUint64(uint64(u))
}

func (u Uint) Size() int {
	return UintValueFieldSize
}

func (u Uint) Type() fields.Type {
	return fields.Uint
}

func (u Uint) IsValid() error {
	return nil
}
