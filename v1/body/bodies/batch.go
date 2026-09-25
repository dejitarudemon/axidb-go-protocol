package bodies

import (
	"sort"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

const (
	FlagsFieldSize          = 1
	RequestsLenFieldSize    = 4
	RequestBodyLenFieldSize = 4

	IsSequentialExecution = 0b1
	InterruptAfterError   = 0b10
	IsOneAnswer           = 0b100
)

type Request struct {
	Number fields.RequestNumber
	Body   body.Body
}

func (r Request) Size() int {
	if r.Body == nil {
		return RequestBodyLenFieldSize + fields.CommandFieldSize + fields.RequestNumberFieldSize
	}
	return RequestBodyLenFieldSize + fields.CommandFieldSize + fields.RequestNumberFieldSize + r.Body.Size()
}

func (r Request) Encode(buf buffer.Appender) {
	r.Number.Encode(buf)

	if r.Body != nil {
		r.Body.Command().Encode(buf)
		buf.AppendUint32(uint32(r.Body.Size()))
		r.Body.Encode(buf)
	} else {
		fields.Command(0).Encode(buf)
		buf.AppendUint32(0)
	}
}

func (r Request) IsValid() error {
	if r.Body == nil {
		return err.NewValidationError(
			"nil Body request",
			"target", "BatchRequest",
			"number", r.Number,
		)
	}

	if err := r.Body.IsValid(); err != nil {
		return err
	}

	switch r.Body.Command() {
	case fields.Read, fields.Write, fields.Delete:
		return nil
	}

	return err.NewValidationError(
		"unexpected fields in batch",
		"target", "BatchRequest",
		"number", r.Number,
		"fields", r.Body.Command(),
	)
}

var _ body.Body = Batch{}

type Batch struct {
	IsSequentialExecution bool
	InterruptAfterError   bool
	IsOneAnswer           bool

	Requests []Request
}

func (b Batch) encodeFlags() uint8 {
	encoded := uint8(0)

	if b.IsSequentialExecution {
		encoded |= IsSequentialExecution
	}

	if b.InterruptAfterError {
		encoded |= InterruptAfterError
	}

	if b.IsOneAnswer {
		encoded |= IsOneAnswer
	}

	return encoded
}

func (b Batch) Size() int {
	size := FlagsFieldSize + RequestsLenFieldSize

	for _, r := range b.Requests {
		size += r.Size()
	}

	return size
}

func (b Batch) Encode(buf buffer.Appender) {
	buf.AppendUint8(b.encodeFlags())
	buf.AppendUint32(uint32(len(b.Requests)))

	for _, r := range b.Requests {
		r.Encode(buf)
	}
}

func (b Batch) Command() fields.Command {
	return fields.Batch
}

func (b Batch) IsValid() error {
	if len(b.Requests) == 0 {
		return err.NewValidationError(
			"no requests",
			"target", b.Command(),
		)
	}
	used := make(map[fields.RequestNumber]int, len(b.Requests))

	for i, r := range b.Requests {
		if err := r.IsValid(); err != nil {
			return err
		}

		if j, ok := used[r.Number]; ok {
			return err.NewValidationError(
				"found at least 2 requests with same numbers",
				"firstIndex", j,
				"secondIndex", i,
				"number", r.Number,
				"target", b.Command(),
			)
		}

		used[r.Number] = i
	}

	return nil
}

func (b *Batch) Sort() {
	sort.Slice(b.Requests, func(i, j int) bool {
		return b.Requests[i].Number < b.Requests[j].Number
	})
}
