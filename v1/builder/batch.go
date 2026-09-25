package builder

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

type BatchRequestsBuilder struct {
	isSequentialExecution bool
	interruptAfterError   bool
	isOneAnswer           bool

	requests []bodies.Request

	cursor fields.RequestNumber
}

func NewBatchRequestsBuilder() *BatchRequestsBuilder {
	return &BatchRequestsBuilder{}
}

func (b *BatchRequestsBuilder) SequentialExecution(v bool) *BatchRequestsBuilder {
	b.isSequentialExecution = v
	return b
}

func (b *BatchRequestsBuilder) InterruptAfterError(v bool) *BatchRequestsBuilder {
	b.interruptAfterError = v
	return b
}

func (b *BatchRequestsBuilder) IsOneAnswer(v bool) *BatchRequestsBuilder {
	b.isOneAnswer = v
	return b
}

func (b *BatchRequestsBuilder) AddRead(key fields.Key) *BatchRequestsBuilder {
	b.requests = append(b.requests, bodies.Request{Number: b.cursor, Body: bodies.Read(key)})
	b.cursor++

	return b
}

func (b *BatchRequestsBuilder) AddWrite(key fields.Key, value value.V) *BatchRequestsBuilder {
	b.requests = append(b.requests, bodies.Request{Number: b.cursor, Body: bodies.Write{Key: key, Value: value}})
	b.cursor++

	return b
}

func (b *BatchRequestsBuilder) AddDelete(key fields.Key) *BatchRequestsBuilder {
	b.requests = append(b.requests, bodies.Request{Number: b.cursor, Body: bodies.Delete(key)})
	b.cursor++

	return b
}
