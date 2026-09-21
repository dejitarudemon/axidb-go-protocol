package bodies

import "github.com/dejitarudemon/axidb-go-protocol/v1/buffer"

const ResultFieldSize = 1

var ResultOK = []byte{0x01}

type simpleOK struct{}

func (s simpleOK) Size() int {
	return ResultFieldSize
}

func (s simpleOK) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
}

func (s simpleOK) IsValid() error {
	return nil
}
