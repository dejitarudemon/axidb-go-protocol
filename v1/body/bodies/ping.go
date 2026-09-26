package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var _ body.Body = Ping{}

// Ping is an empty [fields.Ping] body.
type Ping struct{}

// Size returns the encoded body size in bytes.
func (p Ping) Size() int {
	return 0
}

// Encode writes nothing; a ping body has no payload.
func (p Ping) Encode(buf buffer.Appender) {}

// Command returns [fields.Ping].
func (p Ping) Command() fields.Command {
	return fields.Ping
}

// IsValid reports whether the body satisfies protocol rules.
func (p Ping) IsValid() error {
	return nil
}
