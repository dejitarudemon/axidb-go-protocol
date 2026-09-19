package err

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/google/uuid"
)

const (
	TracebackIDFieldSize = 16
)

type Error interface {
	error

	Size() int
	TracebackID() uuid.UUID
	Encode(buf buffer.Appender)
	Code() Code
}
