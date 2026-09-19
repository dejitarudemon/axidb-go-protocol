package values

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/error/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/types"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = TypedArray{}

const TypedArrayLenFieldSize = 4

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
	for _, elem := range ta.Elems {
		if elem == nil {
			return errs.NewErrorMalformedValue("expected value, got nil")
		}
		if elem.Type() != ta.ElemType {
			return errs.NewErrorMalformedValue(
				fmt.Sprintf("expected %v type, got %v", ta.ElemType, elem.Type()),
			)
		}
		if err := elem.IsValid(); err != nil {
			return err
		}
	}

	return nil
}
