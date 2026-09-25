package builder

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

var testsSpecs_batches = []struct {
	b    *BatchRequestsBuilder
	want bodies.Batch
}{
	{
		b: NewBatchRequestsBuilder().
			InterruptAfterError(true).
			AddRead(fields.Key("key")).
			AddWrite(fields.Key("another-key"), values.String("data")),

		want: bodies.Batch{
			IsSequentialExecution: true,
			InterruptAfterError:   false,
			IsOneAnswer:           true,

			Requests: []bodies.Request{
				{
					Number: 0,
					Body:   bodies.Read("key"),
				},
				{
					Number: 1,
					Body: bodies.Write{
						Key:   fields.Key("another-key"),
						Value: values.String("data"),
					},
				},
			},
		},
	},
}

// var testsSpecs_batch_answers = []struct {
// 	b    *BatchResultsBuilder
// 	want bodies.BatchAnswer
// }{
// 	{
// 		b: NewBatchResultsBuilder().
// 			AddError(1, errs.NewErrorInternalErrorWithTracebackID(
// 				nil,
// 				fields.TracebackID{0x38, 0xDD, 0x5C, 0x39, 0x24, 0xDB, 0x45, 0xA9, 0x85, 0x09, 0xFF, 0xE6, 0xC6, 0x5C, 0x7F, 0x42},
// 			)).
// 			AddError(2, errs.NewErrorRequestInterruptedWithTracebackID(
// 				fields.RequestID(1),
// 				fields.TracebackID{0x49, 0x84, 0x98, 0xC5, 0xCF, 0x19, 0x40, 0xC9, 0x95, 0x38, 0x15, 0x5C, 0x27, 0xA3, 0xCF, 0xF1},
// 			)),
// 		want: bodies.BatchAnswer([]bodies.Result{
// 			{
// 				Number: 1,
// 				Body: bodies.ErrorAnswer{
// 					Err: errs.NewErrorInternalErrorWithTracebackID(
// 						nil,
// 						fields.TracebackID{0x38, 0xDD, 0x5C, 0x39, 0x24, 0xDB, 0x45, 0xA9, 0x85, 0x09, 0xFF, 0xE6, 0xC6, 0x5C, 0x7F, 0x42},
// 					),
// 				},
// 			},
// 			{
// 				Number: 2,
// 				Body: bodies.ErrorAnswer{
// 					Err: errs.NewErrorRequestInterruptedWithTracebackID(
// 						fields.RequestID(1),
// 						fields.TracebackID{0x49, 0x84, 0x98, 0xC5, 0xCF, 0x19, 0x40, 0xC9, 0x95, 0x38, 0x15, 0x5C, 0x27, 0xA3, 0xCF, 0xF1},
// 					),
// 				},
// 			},
// 		},
// 		),
// 	},
// 	{
// 		b: NewBatchResultsBuilder().
// 			AddRead(1, values.String("message")).
// 			AddWrite(2),
// 		want: bodies.BatchAnswer([]bodies.Result{
// 			{
// 				Number: 1,
// 				Body: bodies.ReadAnswer{
// 					Value: values.String("message"),
// 				},
// 			},
// 			{
// 				Number: 2,
// 				Body:   bodies.WriteAnswer{},
// 			},
// 		},
// 		),
// 	},
// }

func compareValues(t *testing.T, v1, v2 value.V) {
	switch want := v2.(type) {
	case values.Bytes:
		got, ok := v1.(values.Bytes)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", v1, want)
			return
		}
		if !bytes.Equal(want, got) {
			t.Fatalf("Compare values.Bytes: got %v, want %v", got, want)
		}
	case values.JSON:
		got, ok := v1.(values.JSON)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", v1, want)
			return
		}
		if !bytes.Equal(want, got) {
			t.Fatalf("Compare values.JSON: got %v, want %v", got, want)
		}
	case values.Int:
		got, ok := v1.(values.Int)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", v1, want)
			return
		}
		if want != got {
			t.Fatalf("Compare values.Int: got %v, want %v", got, want)
		}
	case values.Uint:
		got, ok := v1.(values.Uint)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", v1, want)
			return
		}
		if want != got {
			t.Fatalf("Compare values.Uint: got %v, want %v", got, want)
		}
	case values.Float:
		got, ok := v1.(values.Float)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", v1, want)
			return
		}
		if want != got {
			t.Fatalf("Compare values.Float: got %v, want %v", got, want)
		}
	case values.String:
		got, ok := v1.(values.String)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", v1, want)
			return
		}
		if !strings.EqualFold(string(got), string(want)) {
			t.Fatalf("Compare values.String: got %v, want %v", got, want)
		}
	case values.TypedArray:
		got, ok := v1.(values.TypedArray)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", v1, want)
			return
		}
		if want.ElemType != got.ElemType {
			t.Fatalf("Compare values.TypedArray: mismatched elem type: got %v, want %v", got.ElemType, want.ElemType)
		}
		if len(want.Elems) != len(got.Elems) {
			t.Fatalf("Compare values.TypedArray: mismatched lens: got %v, want %v", len(got.Elems), len(want.Elems))
		}

		for i := range want.Elems {
			compareValues(t, got.Elems[i], want.Elems[i])
		}
	case values.UntypedArray:
		got, ok := v1.(values.UntypedArray)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", v1, want)
			return
		}
		if len(want) != len(got) {
			t.Fatalf("Compare values.TypedArray: mismatched lens: got %v, want %v", len(got), len(want))
		}

		for i := range want {
			compareValues(t, got[i], want[i])
		}
	}
}

func compareBodies(t *testing.T, b1, b2 body.Body) {

	switch want := b2.(type) {
	case bodies.Handshake:
		got, ok := b1.(bodies.Handshake)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}

		if got.Login != want.Login {
			t.Errorf("Compare Logins: got %v and want %v", got.Login, want.Login)
		}

		if !bytes.Equal(got.Hash[:], want.Hash[:]) {
			t.Errorf("Compare Hashes: got % X and want % X", got.Hash, want.Hash)
		}

		if !slices.Equal(got.Compressions, want.Compressions) {
			t.Errorf("Compare Compressions: got %v and want %v", got.Compressions, want.Compressions)
		}
	case bodies.HandshakeAnswer:
		got, ok := b1.(bodies.HandshakeAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}

		if !slices.Equal(got.Compressions, want.Compressions) {
			t.Errorf("Compare Compressions: got %v and want %v", got.Compressions, want.Compressions)
		}
	case bodies.ErrorAnswer:
		got, ok := b1.(bodies.ErrorAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}

		if want.Err.TracebackID() != got.Err.TracebackID() {
			t.Errorf("Compare Traceback ID: got % X and want % X", got.Err.TracebackID(), want.Err.TracebackID())
		}
	case bodies.Read:
		got, ok := b1.(bodies.Read)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}

		if !bytes.Equal([]byte(got), []byte(want)) {
			t.Errorf("Compare Keys: got % X and want % X", got, want)
		}

	case bodies.ReadAnswer:
		got, ok := b1.(bodies.ReadAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}

		compareValues(t, got.Value, want.Value)

	case bodies.Write:
		got, ok := b1.(bodies.Write)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}

		if !bytes.Equal([]byte(got.Key), []byte(want.Key)) {
			t.Errorf("Compare Keys: got % X and want % X", got.Key, want.Key)
		}

		compareValues(t, got.Value, want.Value)
	case bodies.WriteAnswer:
		_, ok := b1.(bodies.WriteAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}
	case bodies.Ping:
		_, ok := b1.(bodies.Ping)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}
	case bodies.PingAnswer:
		_, ok := b1.(bodies.PingAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}
	case bodies.Delete:
		got, ok := b1.(bodies.Delete)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}

		if !bytes.Equal([]byte(got), []byte(want)) {
			t.Errorf("Compare Keys: got % X and want % X", got, want)
		}
	case bodies.DeleteAnswer:
		_, ok := b1.(bodies.DeleteAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}

	case bodies.Batch:
		got, ok := b1.(bodies.Batch)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}

		if got.InterruptAfterError != want.InterruptAfterError {
			t.Errorf("Compare InterruptAfterError: got %v and want %v", got.InterruptAfterError, want.InterruptAfterError)
		}

		if got.IsSequentialExecution != want.IsSequentialExecution {
			t.Errorf("Compare IsSequentialExecution: got %v and want %v", got.IsSequentialExecution, want.IsSequentialExecution)
		}

		if got.IsOneAnswer != want.IsOneAnswer {
			t.Errorf("Compare IsOneAnswer: got %v and want %v", got.IsOneAnswer, want.IsOneAnswer)
		}

		if len(got.Requests) != len(want.Requests) {
			t.Fatalf("Compare Batch: mismatched lens: got %v, want %v", len(got.Requests), len(want.Requests))
		}

		for i := range got.Requests {
			if got.Requests[i].Number != want.Requests[i].Number {
				t.Errorf("Compare Batch Request Numbers: got %v, want %v", got.Requests[i].Number, want.Requests[i].Number)
			}

			compareBodies(t, got.Requests[i].Body, want.Requests[i].Body)
		}
	case bodies.BatchAnswer:
		got, ok := b1.(bodies.BatchAnswer)
		if !ok {
			t.Errorf("Mismatched typed between got %T and want %T", b1, want)
			return
		}

		if len(got) != len(want) {
			t.Fatalf("Compare BatchAnswer: mismatched lens: got %v, want %v", len(got), len(want))
		}

		for i := range got {
			if got[i].Number != want[i].Number {
				t.Errorf("Compare Result Numbers: got %v, want %v", got[i].Number, want[i].Number)
			}

			compareBodies(t, got[i].Body, want[i].Body)
		}
	}
}

func TestBatchRequest_Equal_BySpecs(t *testing.T) {
	for i, tt := range testsSpecs_batches {
		t.Run(
			fmt.Sprintf("TestBatchRequest_Equal_BySpecs %v", i),
			func(t *testing.T) {
				if tt.b == nil {
					t.Fatal("expected BatchRequestsBuilder got nil")
					return
				}

				if len(tt.b.requests) != len(tt.want.Requests) {
					t.Fatalf("mismatched lens beetween got %v and want %v", len(tt.b.requests), len(tt.want.Requests))
					return
				}

				for i := range len(tt.b.requests) {
					if tt.b.requests[i].Number != tt.want.Requests[i].Number {
						t.Errorf("request number: got %v, want %v", tt.b.requests[i].Number, tt.want.Requests[i].Number)
					}

					compareBodies(t, tt.b.requests[i].Body, tt.want.Requests[i].Body)
				}
			},
		)
	}
}

func TestBatchRequest_Equal_Valid(t *testing.T) {
	tests := []struct {
		b    *BatchRequestsBuilder
		want bodies.Batch
	}{
		{
			NewBatchRequestsBuilder(),
			bodies.Batch{},
		},
		{
			NewBatchRequestsBuilder().
				InterruptAfterError(true),
			bodies.Batch{InterruptAfterError: true},
		},
		{
			NewBatchRequestsBuilder().
				SequentialExecution(true),
			bodies.Batch{IsSequentialExecution: true},
		},
		{
			NewBatchRequestsBuilder().
				OneAnswer(true),
			bodies.Batch{IsOneAnswer: true},
		},
		{
			NewBatchRequestsBuilder().
				AddRead(fields.Key("key")),
			bodies.Batch{Requests: []bodies.Request{
				{Number: 0, Body: bodies.Read("key")},
			}},
		},
		{
			NewBatchRequestsBuilder().
				AddDelete(fields.Key("yek")),
			bodies.Batch{Requests: []bodies.Request{
				{Number: 0, Body: bodies.Delete("yek")},
			}},
		},
		{
			NewBatchRequestsBuilder().
				AddWrite(fields.Key("keyek"), values.String("data")),
			bodies.Batch{
				Requests: []bodies.Request{
					{Number: 0, Body: bodies.Write{Key: fields.Key("keyek"), Value: values.String("data")}},
				}},
		},
		{
			NewBatchRequestsBuilder().
				AddRead(fields.Key("key")).
				AddWrite(fields.Key("keyek"), values.String("data")).
				AddDelete(fields.Key("key")),
			bodies.Batch{
				Requests: []bodies.Request{
					{Number: 0, Body: bodies.Read("key")},
					{Number: 1, Body: bodies.Write{Key: fields.Key("keyek"), Value: values.String("data")}},
					{Number: 2, Body: bodies.Delete("key")},
				}},
		},
		{
			NewBatchRequestsBuilder().
				AddWrite(fields.Key("keyek"), values.String("data")).
				AddDelete(fields.Key("key")).
				AddRead(fields.Key("key")),
			bodies.Batch{
				Requests: []bodies.Request{
					{Number: 0, Body: bodies.Write{Key: fields.Key("keyek"), Value: values.String("data")}},
					{Number: 1, Body: bodies.Delete("key")},
					{Number: 2, Body: bodies.Read("key")},
				}},
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatchRequest_Equal_Valid %v", i),
			func(t *testing.T) {
				if tt.b == nil {
					t.Fatal("expected BatchRequestsBuilder got nil")
					return
				}

				if len(tt.b.requests) != len(tt.want.Requests) {
					t.Fatalf("mismatched lens beetween got %v and want %v", len(tt.b.requests), len(tt.want.Requests))
					return
				}

				for i := range len(tt.b.requests) {
					if tt.b.requests[i].Number != tt.want.Requests[i].Number {
						t.Errorf("request number: got %v, want %v", tt.b.requests[i].Number, tt.want.Requests[i].Number)
					}

					compareBodies(t, tt.b.requests[i].Body, tt.want.Requests[i].Body)
				}
			},
		)
	}
}
