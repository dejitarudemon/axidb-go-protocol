package bodies

import (
	"sort"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

const (
	// FlagsFieldSize is the encoded size in bytes of the batch option flags.
	FlagsFieldSize = 1
	// RequestsLenFieldSize is the encoded size in bytes of a request or result count.
	RequestsLenFieldSize = 4
	// RequestBodyLenFieldSize is the encoded size in bytes of a nested body length prefix.
	RequestBodyLenFieldSize = 4

	// IsSequentialExecution is the flag bit that runs batch requests in order.
	IsSequentialExecution = 0b1
	// InterruptAfterError is the flag bit that stops the batch after the first error.
	InterruptAfterError = 0b10
	// IsOneAnswer is the flag bit that replies with a single combined answer.
	IsOneAnswer = 0b100
)

// Request is one numbered command inside a [Batch].
type Request struct {
	// Number identifies the request within the batch.
	Number fields.RequestNumber
	// Body is the nested command payload.
	Body body.Body
}

// Size returns the encoded request size in bytes.
func (r Request) Size() int {
	if r.Body == nil {
		return RequestBodyLenFieldSize + fields.CommandFieldSize + fields.RequestNumberFieldSize
	}
	return RequestBodyLenFieldSize + fields.CommandFieldSize + fields.RequestNumberFieldSize + r.Body.Size()
}

// Encode writes the request number, command, length, and body into buf.
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

// IsValid reports whether the request satisfies protocol rules.
// Body must be non-nil, valid, and a read, write, or delete command.
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

// Batch is a [fields.Batch] body holding several numbered requests and execution flags.
type Batch struct {
	// IsSequentialExecution runs requests in number order when true.
	IsSequentialExecution bool
	// InterruptAfterError stops the batch after the first failed request when true.
	InterruptAfterError bool
	// IsOneAnswer asks for a single combined answer when true.
	IsOneAnswer bool

	// Requests holds the nested commands.
	Requests []Request
}

// encodeFlags packs the batch option bits.
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

// Size returns the encoded body size in bytes.
func (b Batch) Size() int {
	size := FlagsFieldSize + RequestsLenFieldSize

	for _, r := range b.Requests {
		size += r.Size()
	}

	return size
}

// Encode writes the wire encoding of the body into buf.
func (b Batch) Encode(buf buffer.Appender) {
	buf.AppendUint8(b.encodeFlags())
	buf.AppendUint32(uint32(len(b.Requests)))

	for _, r := range b.Requests {
		r.Encode(buf)
	}
}

// Command returns [fields.Batch].
func (b Batch) Command() fields.Command {
	return fields.Batch
}

// IsValid reports whether the body satisfies protocol rules.
// Requests must be non-empty, uniquely numbered, and each must pass its own IsValid check.
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

// Sort orders Requests by ascending request number.
func (b *Batch) Sort() {
	sort.Slice(b.Requests, func(i, j int) bool {
		return b.Requests[i].Number < b.Requests[j].Number
	})
}
