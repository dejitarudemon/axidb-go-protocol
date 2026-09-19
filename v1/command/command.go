package command

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/error/errs"
)

type Code uint8

const FieldSize = 1

const (
	Handshake Code = iota
	Answer
	Read
	Write
	Delete
	Batch
	Ping
)

func (c Code) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(c))
}

func (c Code) String() string {
	switch c {
	case Handshake:
		return "Handshake"
	case Answer:
		return "Answer"
	case Read:
		return "Read"
	case Write:
		return "Write"
	case Delete:
		return "Delete"
	case Batch:
		return "Batch"
	case Ping:
		return "Ping"
	}

	return fmt.Sprintf("Unknown (%d)", c)
}

func (c Code) IsValid() error {
	if c > Ping {
		return errs.NewErrorUnsupportedCommand(uint8(c))
	}

	return nil
}

func (c Code) Size() int {
	return FieldSize
}
