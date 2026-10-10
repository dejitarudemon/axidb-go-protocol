package builder

import (
	"github.com/dejitarudemon/ignicula-wire/v1/body/bodies"
	"github.com/dejitarudemon/ignicula-wire/v1/fields"
	"github.com/dejitarudemon/ignicula-wire/v1/value"
)

// BatchRequestsBuilder assembles numbered read, write, and delete requests for a batch frame.
// Request numbers start at zero and increase by one for each added request.
type BatchRequestsBuilder struct {
	isSequentialExecution bool
	interruptAfterError   bool
	isOneAnswer           bool

	requests []bodies.Request

	cursor fields.RequestNumber
}

// NewBatchRequestsBuilder returns an empty batch request builder.
func NewBatchRequestsBuilder() *BatchRequestsBuilder {
	return &BatchRequestsBuilder{}
}

// SequentialExecution sets whether requests run in number order and returns the builder.
func (b *BatchRequestsBuilder) SequentialExecution(v bool) *BatchRequestsBuilder {
	b.isSequentialExecution = v
	return b
}

// InterruptAfterError sets whether the batch stops after the first error and returns the builder.
func (b *BatchRequestsBuilder) InterruptAfterError(v bool) *BatchRequestsBuilder {
	b.interruptAfterError = v
	return b
}

// OneAnswer sets whether the batch asks for a single combined answer and returns the builder.
func (b *BatchRequestsBuilder) OneAnswer(v bool) *BatchRequestsBuilder {
	b.isOneAnswer = v
	return b
}

// AddRead appends a read request for key and returns the builder.
func (b *BatchRequestsBuilder) AddRead(key fields.Key) *BatchRequestsBuilder {
	b.requests = append(b.requests, bodies.Request{Number: b.cursor, Body: bodies.Read(key)})
	b.cursor++

	return b
}

// AddWrite appends a write request storing value under key and returns the builder.
func (b *BatchRequestsBuilder) AddWrite(key fields.Key, value value.V) *BatchRequestsBuilder {
	b.requests = append(b.requests, bodies.Request{Number: b.cursor, Body: bodies.Write{Key: key, Value: value}})
	b.cursor++

	return b
}

// AddDelete appends a delete request for key and returns the builder.
func (b *BatchRequestsBuilder) AddDelete(key fields.Key) *BatchRequestsBuilder {
	b.requests = append(b.requests, bodies.Request{Number: b.cursor, Body: bodies.Delete(key)})
	b.cursor++

	return b
}

// Extend appends other's requests, numbering them from the receiver cursor, and returns the builder.
// Receiver flags are kept; other's flags are ignored.
func (b *BatchRequestsBuilder) Extend(other *BatchRequestsBuilder) *BatchRequestsBuilder {
	for _, request := range other.requests {
		r := bodies.Request{Number: b.cursor, Body: request.Body}
		b.requests = append(b.requests, r)

		b.cursor++
	}

	return b
}
