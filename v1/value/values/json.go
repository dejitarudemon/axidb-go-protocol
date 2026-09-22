package values

import (
	"encoding/json"
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

var _ value.V = JSON{}

const (
	JSONLenFieldSize       = 4
	FirstSymbolsToShowJSON = 32
)

/*
type JSON представляет собой JSON-документ
из спецификации протокола v1.
*/
type JSON []byte

func (j JSON) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(j)))
	buf.Append(j)
}

func (j JSON) Size() int {
	return JSONLenFieldSize + len(j)
}

func (j JSON) Type() fields.Type {
	return fields.JSON
}

func (j JSON) IsValid() error {
	if ok := json.Valid(j); !ok {
		if len(j) <= FirstSymbolsToShowJSON {
			return err.NewValidationError(
				"invalid JSON format",
				"value", j,
			)
		}

		return err.NewValidationError(
			"invalid JSON format",
			"value", fmt.Sprintf("%q ... and %v symbols", j[:FirstSymbolsToShowJSON], len(j)-FirstSymbolsToShowJSON),
		)

	}

	return nil
}
