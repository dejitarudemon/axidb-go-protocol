package builder

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func TestBatchResultsBuilder(t *testing.T) {
	internalErr := errs.NewErrorInternalErrorWithTracebackID(nil, fields.TracebackID{15: 0x01})
	interrupted := errs.NewErrorRequestInterruptedWithTracebackID(1, fields.TracebackID{15: 0x02})

	tests := []struct {
		name string
		b    *BatchResultsBuilder
		want bodies.BatchAnswer
	}{
		{"empty", NewBatchResultsBuilder(), bodies.BatchAnswer{}},
		{"read", NewBatchResultsBuilder().AddRead(0, values.String("data")), bodies.BatchAnswer{{Number: 0, Body: bodies.ReadAnswer{Value: values.String("data")}}}},
		{"delete", NewBatchResultsBuilder().AddDelete(1), bodies.BatchAnswer{{Number: 1, Body: bodies.DeleteAnswer{}}}},
		{"write", NewBatchResultsBuilder().AddWrite(2), bodies.BatchAnswer{{Number: 2, Body: bodies.WriteAnswer{}}}},
		{"error", NewBatchResultsBuilder().AddError(3, internalErr), bodies.BatchAnswer{{Number: 3, Body: bodies.ErrorAnswer{Err: internalErr}}}},
		{
			"all kinds",
			NewBatchResultsBuilder().
				AddRead(0, values.String("data")).
				AddWrite(1).
				AddDelete(2).
				AddError(3, internalErr),
			bodies.BatchAnswer{
				{Number: 0, Body: bodies.ReadAnswer{Value: values.String("data")}},
				{Number: 1, Body: bodies.WriteAnswer{}},
				{Number: 2, Body: bodies.DeleteAnswer{}},
				{Number: 3, Body: bodies.ErrorAnswer{Err: internalErr}},
			},
		},
		{
			"numbers are kept as given",
			NewBatchResultsBuilder().
				AddWrite(3).
				AddDelete(5).
				AddError(0, interrupted).
				AddRead(10, values.String("qwerty")),
			bodies.BatchAnswer{
				{Number: 3, Body: bodies.WriteAnswer{}},
				{Number: 5, Body: bodies.DeleteAnswer{}},
				{Number: 0, Body: bodies.ErrorAnswer{Err: interrupted}},
				{Number: 10, Body: bodies.ReadAnswer{Value: values.String("qwerty")}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.AssertSameEncoding(t, bodies.BatchAnswer(tt.b.results), tt.want)
		})
	}
}
