package builder

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func TestBatchRequestsBuilder(t *testing.T) {
	tests := []struct {
		name string
		b    *BatchRequestsBuilder
		want bodies.Batch
	}{
		{"empty", NewBatchRequestsBuilder(), bodies.Batch{}},
		{"interrupt after error", NewBatchRequestsBuilder().InterruptAfterError(true), bodies.Batch{InterruptAfterError: true}},
		{"sequential execution", NewBatchRequestsBuilder().SequentialExecution(true), bodies.Batch{IsSequentialExecution: true}},
		{"one answer", NewBatchRequestsBuilder().OneAnswer(true), bodies.Batch{IsOneAnswer: true}},
		{"flag reset", NewBatchRequestsBuilder().OneAnswer(true).OneAnswer(false), bodies.Batch{}},
		{
			"read",
			NewBatchRequestsBuilder().AddRead(fields.Key("key")),
			bodies.Batch{Requests: []bodies.Request{{Number: 0, Body: bodies.Read("key")}}},
		},
		{
			"delete",
			NewBatchRequestsBuilder().AddDelete(fields.Key("yek")),
			bodies.Batch{Requests: []bodies.Request{{Number: 0, Body: bodies.Delete("yek")}}},
		},
		{
			"write",
			NewBatchRequestsBuilder().AddWrite(fields.Key("keyek"), values.String("data")),
			bodies.Batch{Requests: []bodies.Request{{Number: 0, Body: bodies.Write{Key: fields.Key("keyek"), Value: values.String("data")}}}},
		},
		{
			"numbers follow insertion order",
			NewBatchRequestsBuilder().
				AddWrite(fields.Key("keyek"), values.String("data")).
				AddDelete(fields.Key("key")).
				AddRead(fields.Key("key")),
			bodies.Batch{Requests: []bodies.Request{
				{Number: 0, Body: bodies.Write{Key: fields.Key("keyek"), Value: values.String("data")}},
				{Number: 1, Body: bodies.Delete("key")},
				{Number: 2, Body: bodies.Read("key")},
			}},
		},
		{
			"flags and requests",
			NewBatchRequestsBuilder().
				SequentialExecution(true).
				OneAnswer(true).
				AddRead(fields.Key("key")).
				AddWrite(fields.Key("another-key"), values.String("data")),
			bodies.Batch{
				IsSequentialExecution: true,
				IsOneAnswer:           true,
				Requests: []bodies.Request{
					{Number: 0, Body: bodies.Read("key")},
					{Number: 1, Body: bodies.Write{Key: fields.Key("another-key"), Value: values.String("data")}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bodies.Batch{
				IsSequentialExecution: tt.b.isSequentialExecution,
				InterruptAfterError:   tt.b.interruptAfterError,
				IsOneAnswer:           tt.b.isOneAnswer,
				Requests:              tt.b.requests,
			}

			testutil.AssertSameEncoding(t, got, tt.want)
		})
	}
}
