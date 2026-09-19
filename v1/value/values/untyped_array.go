package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/error/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/types"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = UntypedArray{}

const UntypedArrayLenFieldSize = 4

type UntypedArray []value.V

func (ua UntypedArray) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(ua)))

	for _, elem := range ua {
		elem.Type().Encode(buf)
		elem.Encode(buf)
	}
}

func (ua UntypedArray) Size() int {
	size := 0

	for _, elem := range ua {
		size += elem.Size() + types.FieldSize
	}

	return UntypedArrayLenFieldSize + types.FieldSize + size
}

func (ua UntypedArray) Type() types.Code {
	return types.UntypedArray
}

func (ua UntypedArray) IsValid() error {
	for _, elem := range ua {
		if elem == nil {
			return errs.NewErrorMalformedValue("expected value, got nil")
		}

		if err := elem.IsValid(); err != nil {
			return err
		}
	}

	return nil
}
