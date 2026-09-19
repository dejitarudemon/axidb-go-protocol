package values

import (
	"math"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/types"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = Float(0)

const FloatValueFieldSize = 8

type Float float64

func (f Float) Encode(buf buffer.Appender) {
	buf.AppendUint64(math.Float64bits(float64(f)))
}

func (f Float) Size() int {
	return FloatValueFieldSize
}

func (f Float) Type() types.Code {
	return types.Float
}

func (f Float) IsValid() error {
	return nil
}
