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
	// JSONLenFieldSize is the encoded size in bytes of the length prefix for [JSON].
	JSONLenFieldSize = 4
	// FirstSymbolsToShowJSON is the number of bytes included in validation error details for large invalid JSON.
	FirstSymbolsToShowJSON = 32
)

// JSON is a JSON document with type code [fields.JSON].
type JSON []byte

// Encode writes the wire encoding of the value into buf.
func (j JSON) Encode(buf buffer.Appender) {
	buf.AppendUint32(uint32(len(j)))
	buf.Append(j)
}

// Size returns the encoded value size in bytes.
func (j JSON) Size() int {
	return JSONLenFieldSize + len(j)
}

// Type returns [fields.JSON].
func (j JSON) Type() fields.Type {
	return fields.JSON
}

// IsValid reports whether the value satisfies type-specific rules.
// The payload must be valid JSON according to [encoding/json.Valid].
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
