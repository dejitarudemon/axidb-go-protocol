package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/types"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = TypedArray{}

const TypedArrayLenFieldSize = 4

/*
type TypedArray представляет собой типизированную последовательность
элементов из спецификации протокола v1.
*/
type TypedArray struct {
	ElemType types.Code
	Elems    []value.V
}

func (ta TypedArray) realLen() int {
	realLen := 0
	for _, elem := range ta.Elems {
		if elem != nil {
			realLen += 1
		}
	}

	return realLen
}

func (ta TypedArray) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(ta.realLen()))
	ta.ElemType.Encode(buf)

	for _, elem := range ta.Elems {
		if elem != nil {
			elem.Encode(buf)
		}
	}
}

func (ta TypedArray) Size() int {
	size := 0

	for _, elem := range ta.Elems {
		if elem != nil {
			size += elem.Size()
		}
	}

	return TypedArrayLenFieldSize + types.FieldSize + size
}

func (ta TypedArray) Type() types.Code {
	return types.TypedArray
}

func (ta TypedArray) IsValid() error {
	for i, elem := range ta.Elems {
		if elem == nil {
			return err.NewValidationError(
				"nil elem in TypedArray",
				"index", i,
			)
		}
		if elem.Type() != ta.ElemType {
			return err.NewValidationError(
				"wrong elem's type in TypedArray",
				"expected", ta.ElemType,
				"got", elem.Type(),
				"index", i,
			)
		}
		if err := elem.IsValid(); err != nil {
			return err
		}
	}

	return nil
}
