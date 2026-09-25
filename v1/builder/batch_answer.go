package builder

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

type BatchResultsBuilder struct {
	results []bodies.Result
}

func (b *BatchResultsBuilder) AddRead(number fields.RequestNumber, value value.V) *BatchResultsBuilder {
	b.results = append(b.results, bodies.Result{Number: number, Body: bodies.ReadAnswer{Value: value}})

	return b
}

func (b *BatchResultsBuilder) AddWrite(number fields.RequestNumber) *BatchResultsBuilder {
	b.results = append(b.results, bodies.Result{Number: number, Body: bodies.WriteAnswer{}})

	return b
}

func (b *BatchResultsBuilder) AddDelete(number fields.RequestNumber) *BatchResultsBuilder {
	b.results = append(b.results, bodies.Result{Number: number, Body: bodies.DeleteAnswer{}})

	return b
}

func (b *BatchResultsBuilder) AddError(number fields.RequestNumber, err err.ProtocolError) *BatchResultsBuilder {
	b.results = append(b.results, bodies.Result{Number: number, Body: bodies.ErrorAnswer{Err: err}})

	return b
}
