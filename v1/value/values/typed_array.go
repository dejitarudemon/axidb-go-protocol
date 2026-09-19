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
	elemType types.Code
	elems    []value.V
}

func (ta TypedArray) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(ta.elems)))
	ta.elemType.Encode(buf)

	for _, elem := range ta.elems {
		elem.Encode(buf)
	}
}

func (ta TypedArray) Size() int {
	size := 0

	for _, elem := range ta.elems {
		size += elem.Size()
	}

	return TypedArrayLenFieldSize + types.FieldSize + size
}

func (ta TypedArray) Type() types.Code {
	return types.TypedArray
}

func (ta TypedArray) IsValid() error {
	for _, elem := range ta.elems {
		if elem == nil {
			return errs.NewErrorMalformedValue("expected value, got nil")
		}
		if elem.Type() != ta.elemType {
			return errs.NewErrorMalformedValue(
				fmt.Sprintf("expected %v type, got %v", ta.elemType, elem.Type()),
			)
		}
		if err := elem.IsValid(); err != nil {
			return err
		}
	}

	return nil
}
