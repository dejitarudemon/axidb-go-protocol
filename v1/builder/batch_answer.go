package builder

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

// BatchResultsBuilder assembles numbered answers for a batch answer frame.
type BatchResultsBuilder struct {
	results []bodies.Result
}

// NewBatchResultsBuilder returns an empty batch result builder.
func NewBatchResultsBuilder() *BatchResultsBuilder {
	return &BatchResultsBuilder{}
}

// AddRead appends a successful read answer for number carrying value and returns the builder.
func (b *BatchResultsBuilder) AddRead(number fields.RequestNumber, value value.V) *BatchResultsBuilder {
	b.results = append(b.results, bodies.Result{Number: number, Body: bodies.ReadAnswer{Value: value}})

	return b
}

// AddWrite appends a successful write answer for number and returns the builder.
func (b *BatchResultsBuilder) AddWrite(number fields.RequestNumber) *BatchResultsBuilder {
	b.results = append(b.results, bodies.Result{Number: number, Body: bodies.WriteAnswer{}})

	return b
}

// AddDelete appends a successful delete answer for number and returns the builder.
func (b *BatchResultsBuilder) AddDelete(number fields.RequestNumber) *BatchResultsBuilder {
	b.results = append(b.results, bodies.Result{Number: number, Body: bodies.DeleteAnswer{}})

	return b
}

// AddError appends an error answer for number and returns the builder.
func (b *BatchResultsBuilder) AddError(number fields.RequestNumber, err err.ProtocolError) *BatchResultsBuilder {
	b.results = append(b.results, bodies.Result{Number: number, Body: bodies.ErrorAnswer{Err: err}})

	return b
}
