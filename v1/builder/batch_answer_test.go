package builder

import (
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

var testsSpecs_batch_answers = []struct {
	b    *BatchResultsBuilder
	want bodies.BatchAnswer
}{
	{
		b: NewBatchResultsBuilder().
			AddError(1, errs.NewErrorInternalErrorWithTracebackID(
				nil,
				fields.TracebackID{0x38, 0xDD, 0x5C, 0x39, 0x24, 0xDB, 0x45, 0xA9, 0x85, 0x09, 0xFF, 0xE6, 0xC6, 0x5C, 0x7F, 0x42},
			)).
			AddError(2, errs.NewErrorRequestInterruptedWithTracebackID(
				fields.RequestID(1),
				fields.TracebackID{0x49, 0x84, 0x98, 0xC5, 0xCF, 0x19, 0x40, 0xC9, 0x95, 0x38, 0x15, 0x5C, 0x27, 0xA3, 0xCF, 0xF1},
			)),
		want: bodies.BatchAnswer([]bodies.Result{
			{
				Number: 1,
				Body: bodies.ErrorAnswer{
					Err: errs.NewErrorInternalErrorWithTracebackID(
						nil,
						fields.TracebackID{0x38, 0xDD, 0x5C, 0x39, 0x24, 0xDB, 0x45, 0xA9, 0x85, 0x09, 0xFF, 0xE6, 0xC6, 0x5C, 0x7F, 0x42},
					),
				},
			},
			{
				Number: 2,
				Body: bodies.ErrorAnswer{
					Err: errs.NewErrorRequestInterruptedWithTracebackID(
						fields.RequestID(1),
						fields.TracebackID{0x49, 0x84, 0x98, 0xC5, 0xCF, 0x19, 0x40, 0xC9, 0x95, 0x38, 0x15, 0x5C, 0x27, 0xA3, 0xCF, 0xF1},
					),
				},
			},
		},
		),
	},
	{
		b: NewBatchResultsBuilder().
			AddRead(1, values.String("message")).
			AddWrite(2),
		want: bodies.BatchAnswer([]bodies.Result{
			{
				Number: 1,
				Body: bodies.ReadAnswer{
					Value: values.String("message"),
				},
			},
			{
				Number: 2,
				Body:   bodies.WriteAnswer{},
			},
		},
		),
	},
}

func TestBatchResult_Equal_BySpecs(t *testing.T) {
	for i, tt := range testsSpecs_batch_answers {
		t.Run(
			fmt.Sprintf("TestBatchResult_Equal_BySpecs %v", i),
			func(t *testing.T) {
				if tt.b == nil {
					t.Fatal("expected BatchResultsBuilder got nil")
					return
				}

				if len(tt.b.results) != len(tt.want) {
					t.Fatalf("mismatched lens beetween got %v and want %v", len(tt.b.results), len(tt.want))
					return
				}

				for i := range len(tt.b.results) {
					if tt.b.results[i].Number != tt.want[i].Number {
						t.Errorf("request number: got %v, want %v", tt.b.results[i].Number, tt.want[i].Number)
					}

					compareBodies(t, tt.b.results[i].Body, tt.want[i].Body)
				}
			},
		)
	}
}

func TestBatchResult_Equal_Valid(t *testing.T) {
	tests := []struct {
		b    *BatchResultsBuilder
		want bodies.Batch
	}{
		{
			NewBatchResultsBuilder(),
			bodies.Batch{},
		},
		{
			NewBatchResultsBuilder().
				AddRead(0, values.String("data")),
			bodies.Batch{Requests: []bodies.Request{
				{Number: 0, Body: bodies.ReadAnswer{Value: values.String("data")}},
			}},
		},
		{
			NewBatchResultsBuilder().
				AddDelete(1),
			bodies.Batch{Requests: []bodies.Request{
				{Number: 1, Body: bodies.DeleteAnswer{}},
			}},
		},
		{
			NewBatchResultsBuilder().
				AddWrite(2),
			bodies.Batch{
				Requests: []bodies.Request{
					{Number: 2, Body: bodies.WriteAnswer{}},
				}},
		},
		{
			NewBatchResultsBuilder().
				AddRead(0, values.String("data")).
				AddWrite(1).
				AddDelete(2).
				AddError(3, errs.NewErrorInternalErrorWithTracebackID(nil, fields.TracebackID{})),
			bodies.Batch{
				Requests: []bodies.Request{
					{Number: 0, Body: bodies.ReadAnswer{Value: values.String("data")}},
					{Number: 1, Body: bodies.WriteAnswer{}},
					{Number: 2, Body: bodies.DeleteAnswer{}},
					{Number: 3, Body: bodies.ErrorAnswer{Err: errs.NewErrorInternalErrorWithTracebackID(nil, fields.TracebackID{})}},
				}},
		},
		{
			NewBatchResultsBuilder().
				AddWrite(3).
				AddDelete(5).
				AddError(0, errs.NewErrorInternalErrorWithTracebackID(nil, fields.TracebackID{})).
				AddRead(10, values.String("qwerty")),
			bodies.Batch{
				Requests: []bodies.Request{
					{Number: 3, Body: bodies.WriteAnswer{}},
					{Number: 5, Body: bodies.DeleteAnswer{}},
					{Number: 0, Body: bodies.ErrorAnswer{Err: errs.NewErrorInternalErrorWithTracebackID(nil, fields.TracebackID{})}},
					{Number: 10, Body: bodies.ReadAnswer{Value: values.String("qwerty")}},
				}},
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatchResult_Equal_Valid %v", i),
			func(t *testing.T) {
				if tt.b == nil {
					t.Fatal("expected BatchResultsBuilder got nil")
					return
				}

				if len(tt.b.results) != len(tt.want.Requests) {
					t.Fatalf("mismatched lens beetween got %v and want %v", len(tt.b.results), len(tt.want.Requests))
					return
				}

				for i := range len(tt.b.results) {
					if tt.b.results[i].Number != tt.want.Requests[i].Number {
						t.Errorf("request number: got %v, want %v", tt.b.results[i].Number, tt.want.Requests[i].Number)
					}

					compareBodies(t, tt.b.results[i].Body, tt.want.Requests[i].Body)
				}
			},
		)
	}
}
