package bodies

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

type Result struct {
	Number uint32
	Body   body.Body
}

func (r Result) Size() int {
	if r.Body == nil {
		return RequestNumberFieldSize + RequestNumberFieldSize
	}
	return RequestsLenFieldSize + RequestNumberFieldSize + r.Body.Size()
}

func (r Result) Encode(buf buffer.Appender) {
	buf.AppendUint32(r.Number)

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

var _ body.Body = BatchAnswer{}

type BatchAnswer struct {
	Results []Result
}

func (b BatchAnswer) Size() int {
	size := RequestsLenFieldSize + ResultFieldSize

	for _, r := range b.Results {
		size += r.Size()
	}

	return size
}

func (b BatchAnswer) Encode(buf buffer.Appender) {
	buf.Append(ResultOK)
	buf.AppendUint32(uint32(len(b.Results)))

	for _, r := range b.Results {
		r.Encode(buf)
	}
}

func (b BatchAnswer) Command() fields.Command {
	return fields.Answer
}

func (b BatchAnswer) IsValid() error {
	if len(b.Results) == 0 {
		return err.NewValidationError(
			"no results",
			"target", b.Command(),
		)
	}
	used := make(map[uint32]int, len(b.Results))

	for i, r := range b.Results {
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
