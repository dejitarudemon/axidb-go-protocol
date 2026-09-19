package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/types"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = Bytes([]byte{})

const BytesLenFieldSize = 4

type Bytes []byte

func (b Bytes) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(b)))
	buf.Append(b)
}

func (b Bytes) Size() int {
	return BytesLenFieldSize + len(b)
}

func (b Bytes) Type() types.Code {
	return types.Bytes
}

func (b Bytes) IsValid() error {
	return nil
}
