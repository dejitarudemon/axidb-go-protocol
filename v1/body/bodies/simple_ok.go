package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

// ResultFieldSize is the encoded size in bytes of the success or failure flag in an answer.
const ResultFieldSize = 1

// ResultOK is the success flag written at the start of a successful answer.
var ResultOK = []byte{0x01}

// simpleOK is the shared size and validation of answers that contain only a success flag and a command code.
type simpleOK struct{}

// Size returns the encoded size of a success flag and a command code.
func (s simpleOK) Size() int {
	return ResultFieldSize + fields.CommandFieldSize
}

// IsValid reports that a success-only answer has no additional payload to check.
func (s simpleOK) IsValid() error {
	return nil
}
