package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/types"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

const (
	KeyLenFieldSize = 4
)

var _ body.Body = Write{}

type Write struct {
	Key   fields.Key
	Value value.V
}

func (w Write) Size() int {
	if w.Value == nil {
		return len(w.Key) + KeyLenFieldSize
	}

	return len(w.Key) + KeyLenFieldSize + w.Value.Size() + types.FieldSize
}

func (w Write) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(w.Key.Size()))
	w.Key.Encode(buf)

	if w.Value != nil {
		w.Value.Type().Encode(buf)
		w.Value.Encode(buf)
	}
}

func (w Write) Command() fields.Command {
	return fields.Write
}

func (w Write) IsValid() error {
	if len(w.Key) == 0 {
		return err.NewValidationError(
			"empty key",
			"target", w.Command(),
		)
	}

	if w.Value == nil {
		return err.NewValidationError(
			"nil value",
			"target", w.Command(),
		)
	}

	return w.Value.IsValid()
}
