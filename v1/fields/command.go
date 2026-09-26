package fields

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// Command is a protocol v1 frame command code.
type Command uint8

// CommandFieldSize is the encoded size in bytes of a [Command].
const CommandFieldSize = 1

// Command codes defined by protocol v1.
const (
	Handshake Command = iota
	Answer
	Read
	Write
	Delete
	Batch
	Ping
)

// Encode writes the wire encoding of the command into buf.
func (c Command) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(c))
}

// String returns a human-readable command name.
func (c Command) String() string {
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

// IsValid reports whether c is a known protocol command (0 through [Ping]).
func (c Command) IsValid() bool {
	return c <= Ping
}

// Size returns the encoded size in bytes.
func (c Command) Size() int {
	return CommandFieldSize
}
