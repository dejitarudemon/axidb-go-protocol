package values

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = UntypedArray{}

const UntypedArrayLenFieldSize = 4

/*
type UntypedArray представляет собой нетипизированную
последовательность элементов из спецификации протокола v1.
*/
type UntypedArray []value.V

func (ua UntypedArray) realLen() int {
	realLen := 0
	for _, elem := range ua {
		if elem != nil {
			realLen += 1
		}
	}

	return realLen
}

func (ua UntypedArray) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(ua.realLen()))

	for _, elem := range ua {
		if elem != nil {
			elem.Type().Encode(buf)
			elem.Encode(buf)
		}
	}
}

func (ua UntypedArray) Size() int {
	size := UntypedArrayLenFieldSize

	for _, elem := range ua {
		if elem != nil {
			size += elem.Size() + elem.Type().Size()
		}
	}

	return size
}

func (ua UntypedArray) Type() fields.Type {
	return fields.UntypedArray
}

func (ua UntypedArray) IsValid() error {
	for i, elem := range ua {
		if elem == nil {
			return err.NewValidationError(
				"nil elem in UntypedArray",
				"index", i,
			)
		}

		if err := elem.IsValid(); err != nil {
			return err
		}
	}

	return nil
}
