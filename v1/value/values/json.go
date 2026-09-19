package values

import (
	"encoding/json"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/error/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/types"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = JSON{}

const JSONLenFieldSize = 4

type JSON []byte

func (j JSON) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(j)))
	buf.Append(j)
}

func (j JSON) Size() int {
	return JSONLenFieldSize + len(j)
}

func (j JSON) Type() types.Code {
	return types.JSON
}

func (j JSON) IsValid() error {
	if ok := json.Valid(j); !ok {
		return errs.NewErrorMalformedValue(
			"invalid JSON format",
		)
	}

	return nil
}
