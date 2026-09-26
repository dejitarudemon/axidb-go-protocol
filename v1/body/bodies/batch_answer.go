package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

// Result is one numbered answer inside a [BatchAnswer].
type Result struct {
	// Number identifies the request this result replies to.
	Number fields.RequestNumber
	// Body is the nested answer.
	Body body.Answer
}

// Size returns the encoded result size in bytes.
func (r Result) Size() int {
	if r.Body == nil {
		return RequestsLenFieldSize + fields.RequestNumberFieldSize
	}
	return RequestsLenFieldSize + fields.RequestNumberFieldSize + r.Body.Size()
}

// Encode writes the request number, answer length, and answer into buf.
func (r Result) Encode(buf buffer.Appender) {
	r.Number.Encode(buf)

	if r.Body != nil {
		buf.AppendUint32(uint32(r.Body.Size()))
		r.Body.Encode(buf)
	} else {
		buf.AppendUint32(0)
	}
}

// IsResponseTo returns the command the nested answer replies to.
// A nil body reports command code 0.
func (r Result) IsResponseTo() fields.Command {
	if r.Body == nil {
		return fields.Command(0)
	}
	return r.Body.IsResponseTo()
}

// IsValid reports whether the result satisfies protocol rules.
// Body must be a non-nil answer whose command is [fields.Answer] and which passes its own IsValid check.
func (r Result) IsValid() error {
	if r.Body == nil {
		return err.NewValidationError(
			"nil Body result",
			"target", "BatchResult",
			"number", r.Number,
		)
	}

	if err := r.Body.IsValid(); err != nil {
		return err
	}

	if r.Body.Command() == fields.Answer {
		return nil
	}

	return err.NewValidationError(
		"not an answer in batch",
		"target", "BatchResult",
		"number", r.Number,
		"fields", r.Body.Command(),
	)
}

var _ body.Answer = BatchAnswer{}

// BatchAnswer is a successful reply to [fields.Batch] listing one result per request.
type BatchAnswer []Result

// Size returns the encoded body size in bytes.
func (b BatchAnswer) Size() int {
	size := RequestsLenFieldSize + ResultFieldSize + b.IsResponseTo().Size()

	for _, r := range b {
		size += r.Size()
	}

	return size
}

// Encode writes the wire encoding of the body into buf.
func (b BatchAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	b.IsResponseTo().Encode(buf)
	buf.AppendUint32(uint32(len(b)))

	for _, r := range b {
		r.Encode(buf)
	}
}

// Command returns [fields.Answer].
func (b BatchAnswer) Command() fields.Command {
	return fields.Answer
}

// IsResponseTo returns [fields.Batch].
func (b BatchAnswer) IsResponseTo() fields.Command {
	return fields.Batch
}

// IsValid reports whether the body satisfies protocol rules.
// Results must be non-empty, uniquely numbered, and each must pass its own IsValid check.
func (b BatchAnswer) IsValid() error {
	if len(b) == 0 {
		return err.NewValidationError(
			"no results",
			"target", b.Command(),
		)
	}
	used := make(map[fields.RequestNumber]int, len(b))

	for i, r := range b {
		if err := r.IsValid(); err != nil {
			return err
		}

		if j, ok := used[r.Number]; ok {
			return err.NewValidationError(
				"found at least 2 results with same numbers",
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
