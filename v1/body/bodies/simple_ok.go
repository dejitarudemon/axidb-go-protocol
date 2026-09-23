package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

const ResultFieldSize = 1

var ResultOK = []byte{0x01}

type simpleOK struct{}

func (s simpleOK) Size() int {
	return ResultFieldSize + fields.CommandFieldSize
}

func (s simpleOK) IsValid() error {
	return nil
}
