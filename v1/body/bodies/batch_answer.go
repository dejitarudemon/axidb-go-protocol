package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

type Result struct {
	Number fields.RequestNumber
	Body   body.Answer
}

func (r Result) Size() int {
	if r.Body == nil {
		return RequestsLenFieldSize + fields.RequestNumberFieldSize
	}
	return RequestsLenFieldSize + fields.RequestNumberFieldSize + r.Body.Size()
}

func (r Result) Encode(buf buffer.Appender) {
	r.Number.Encode(buf)

	if r.Body != nil {
		buf.AppendUint32(uint32(r.Body.Size()))
		r.Body.Encode(buf)
	} else {
		buf.AppendUint32(0)
	}
}

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

type BatchAnswer []Result

func (b BatchAnswer) Size() int {
	size := RequestsLenFieldSize + ResultFieldSize + b.IsResponseTo().Size()

	for _, r := range b {
		size += r.Size()
	}

	return size
}

func (b BatchAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	b.IsResponseTo().Encode(buf)
	buf.AppendUint32(uint32(len(b)))

	for _, r := range b {
		r.Encode(buf)
	}
}

func (b BatchAnswer) Command() fields.Command {
	return fields.Answer
}

func (b BatchAnswer) IsResponseTo() fields.Command {
	return fields.Batch
}

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
