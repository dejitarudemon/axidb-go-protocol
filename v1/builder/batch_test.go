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
			assertBatchRequests(t, tt.b, tt.want)
		})
	}
}

func TestBatchRequestsBuilder_Extend(t *testing.T) {
	tests := []struct {
		name string
		b    *BatchRequestsBuilder
		want bodies.Batch
	}{
		{
			"empty other",
			NewBatchRequestsBuilder().AddRead(fields.Key("key")).Extend(NewBatchRequestsBuilder()),
			bodies.Batch{Requests: []bodies.Request{{Number: 0, Body: bodies.Read("key")}}},
		},
		{
			"into empty",
			NewBatchRequestsBuilder().Extend(
				NewBatchRequestsBuilder().AddRead(fields.Key("a")).AddDelete(fields.Key("b")),
			),
			bodies.Batch{Requests: []bodies.Request{
				{Number: 0, Body: bodies.Read("a")},
				{Number: 1, Body: bodies.Delete("b")},
			}},
		},
		{
			"appends and continues numbers",
			NewBatchRequestsBuilder().
				AddRead(fields.Key("a")).
				Extend(NewBatchRequestsBuilder().AddWrite(fields.Key("b"), values.String("v")).AddDelete(fields.Key("c"))),
			bodies.Batch{Requests: []bodies.Request{
				{Number: 0, Body: bodies.Read("a")},
				{Number: 1, Body: bodies.Write{Key: fields.Key("b"), Value: values.String("v")}},
				{Number: 2, Body: bodies.Delete("c")},
			}},
		},
		{
			"keeps receiver flags",
			NewBatchRequestsBuilder().
				SequentialExecution(true).
				InterruptAfterError(true).
				Extend(NewBatchRequestsBuilder().OneAnswer(true).AddRead(fields.Key("k"))),
			bodies.Batch{
				IsSequentialExecution: true,
				InterruptAfterError:   true,
				Requests:              []bodies.Request{{Number: 0, Body: bodies.Read("k")}},
			},
		},
		{
			"renumbers other from receiver cursor",
			NewBatchRequestsBuilder().
				AddDelete(fields.Key("x")).
				AddRead(fields.Key("y")).
				Extend(NewBatchRequestsBuilder().AddRead(fields.Key("z"))),
			bodies.Batch{Requests: []bodies.Request{
				{Number: 0, Body: bodies.Delete("x")},
				{Number: 1, Body: bodies.Read("y")},
				{Number: 2, Body: bodies.Read("z")},
			}},
		},
		{
			"add after extend continues cursor",
			NewBatchRequestsBuilder().
				AddRead(fields.Key("a")).
				Extend(NewBatchRequestsBuilder().AddDelete(fields.Key("b"))).
				AddWrite(fields.Key("c"), values.String("v")),
			bodies.Batch{Requests: []bodies.Request{
				{Number: 0, Body: bodies.Read("a")},
				{Number: 1, Body: bodies.Delete("b")},
				{Number: 2, Body: bodies.Write{Key: fields.Key("c"), Value: values.String("v")}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBatchRequests(t, tt.b, tt.want)
		})
	}
}

func assertBatchRequests(t *testing.T, b *BatchRequestsBuilder, want bodies.Batch) {
	t.Helper()

	got := bodies.Batch{
		IsSequentialExecution: b.isSequentialExecution,
		InterruptAfterError:   b.interruptAfterError,
		IsOneAnswer:           b.isOneAnswer,
		Requests:              b.requests,
	}

	testutil.AssertSameEncoding(t, got, want)
}
