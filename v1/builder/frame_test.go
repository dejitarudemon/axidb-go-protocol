package builder

import (
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func TestFrameBuilder_Ping(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		r       fields.RequestID
		want    frame.Frame
		wantErr bool
	}{
		{
			0,
			frame.Frame{
				RequestID: 0,
				Body:      bodies.Ping{},
			},
			true,
		},
		{
			1,
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Ping{},
			},
			false,
		},
		{
			math.MaxUint32,
			frame.Frame{
				RequestID: math.MaxUint32,
				Body:      bodies.Ping{},
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_Ping %v", i),
			func(t *testing.T) {
				f, err := fb.NewPing(tt.r)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_Read(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		r       fields.RequestID
		k       fields.Key
		want    frame.Frame
		wantErr bool
	}{
		{
			0,
			fields.Key(""),
			frame.Frame{
				RequestID: 0,
				Body:      bodies.Read(""),
			},
			true,
		},
		{
			0,
			fields.Key("key"),
			frame.Frame{
				RequestID: 0,
				Body:      bodies.Read("key"),
			},
			true,
		},
		{
			1,
			fields.Key(""),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Read(""),
			},
			true,
		},
		{
			1,
			fields.Key("key"),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Read("key"),
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_Read %v", i),
			func(t *testing.T) {
				f, err := fb.NewRead(tt.r, tt.k)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_Delete(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		r       fields.RequestID
		k       fields.Key
		want    frame.Frame
		wantErr bool
	}{
		{
			0,
			fields.Key(""),
			frame.Frame{
				RequestID: 0,
				Body:      bodies.Delete(""),
			},
			true,
		},
		{
			0,
			fields.Key("key"),
			frame.Frame{
				RequestID: 0,
				Body:      bodies.Delete("key"),
			},
			true,
		},
		{
			1,
			fields.Key(""),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Delete(""),
			},
			true,
		},
		{
			1,
			fields.Key("key"),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Delete("key"),
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_Delete %v", i),
			func(t *testing.T) {
				f, err := fb.NewDelete(tt.r, tt.k)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_Write(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		r       fields.RequestID
		k       fields.Key
		v       value.V
		want    frame.Frame
		wantErr bool
	}{
		{
			0,
			fields.Key(""),
			nil,
			frame.Frame{
				RequestID: 0,
				Body:      bodies.Write{},
			},
			true,
		},
		{
			1,
			fields.Key(""),
			nil,
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Write{},
			},
			true,
		},
		{
			1,
			fields.Key("key"),
			nil,
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Write{Key: fields.Key("key")},
			},
			true,
		},
		{
			1,
			fields.Key("key"),
			values.JSON("-"),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Write{Key: fields.Key("key"), Value: values.JSON("-")},
			},
			true,
		},
		{
			1,
			fields.Key("key"),
			values.Int(3),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Write{Key: fields.Key("key"), Value: values.Int(3)},
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_Write %v", i),
			func(t *testing.T) {
				f, err := fb.NewWrite(tt.r, tt.k, tt.v)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_Handshake(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		l       string
		h       [32]byte
		c       []fields.Compression
		want    frame.Frame
		wantErr bool
	}{
		{
			"",
			[32]byte{},
			nil,
			frame.Frame{
				Body: bodies.Handshake{},
			},
			false,
		},
		{
			"user",
			[32]byte{},
			nil,
			frame.Frame{
				Body: bodies.Handshake{Login: "user"},
			},
			false,
		},
		{
			"user",
			[32]byte{},
			nil,
			frame.Frame{
				Body: bodies.Handshake{Login: "user"},
			},
			false,
		},
		{
			"",
			[32]byte{},
			slices.Repeat([]fields.Compression{fields.Zstd}, int(bodies.MaxCompressionsPerOneHandshake)+1),
			frame.Frame{
				RequestID: 0,
				Body:      bodies.Handshake{Compressions: []fields.Compression{fields.Zstd}},
			},
			false,
		},
		{
			"",
			[32]byte{},
			slices.Repeat([]fields.Compression{fields.Zstd}, int(bodies.MaxCompressionsPerOneHandshake)+1),
			frame.Frame{
				RequestID: 0,
				Body:      bodies.Handshake{Compressions: []fields.Compression{fields.Zstd}},
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_Handshake %v", i),
			func(t *testing.T) {
				f, err := fb.NewHandshake(tt.l, tt.h, tt.c)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_WriteAnswer(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		r       fields.RequestID
		want    frame.Frame
		wantErr bool
	}{
		{
			0,
			frame.Frame{
				RequestID: 0,
				Body:      bodies.WriteAnswer{},
			},
			true,
		},
		{
			1,
			frame.Frame{
				RequestID: 1,
				Body:      bodies.WriteAnswer{},
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_WriteAnswer %v", i),
			func(t *testing.T) {
				f, err := fb.NewWriteAnswer(tt.r)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_DeleteAnswer(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		r       fields.RequestID
		want    frame.Frame
		wantErr bool
	}{
		{
			0,
			frame.Frame{
				RequestID: 0,
				Body:      bodies.DeleteAnswer{},
			},
			true,
		},
		{
			1,
			frame.Frame{
				RequestID: 1,
				Body:      bodies.DeleteAnswer{},
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_DeleteAnswer %v", i),
			func(t *testing.T) {
				f, err := fb.NewDeleteAnswer(tt.r)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_PingAnswer(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		r       fields.RequestID
		want    frame.Frame
		wantErr bool
	}{
		{
			0,
			frame.Frame{
				RequestID: 0,
				Body:      bodies.PingAnswer{},
			},
			true,
		},
		{
			1,
			frame.Frame{
				RequestID: 1,
				Body:      bodies.PingAnswer{},
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_PingAnswer %v", i),
			func(t *testing.T) {
				f, err := fb.NewPingAnswer(tt.r)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_ReadAnswer(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		r       fields.RequestID
		v       value.V
		want    frame.Frame
		wantErr bool
	}{
		{
			0,
			nil,
			frame.Frame{
				RequestID: 0,
				Body:      bodies.ReadAnswer{},
			},
			true,
		},
		{
			1,
			nil,
			frame.Frame{
				RequestID: 1,
				Body:      bodies.ReadAnswer{},
			},
			true,
		},
		{
			1,
			values.JSON("-"),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.ReadAnswer{Value: values.JSON("-")},
			},
			true,
		},
		{
			1,
			values.Int(3),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.ReadAnswer{Value: values.Int(3)},
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_ReadAnswer %v", i),
			func(t *testing.T) {
				f, err := fb.NewReadAnswer(tt.r, tt.v)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_HandshakeAnswer(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		c       []fields.Compression
		want    frame.Frame
		wantErr bool
	}{
		{
			nil,
			frame.Frame{
				Body: bodies.HandshakeAnswer{},
			},
			false,
		},
		{
			slices.Repeat([]fields.Compression{fields.Zstd}, int(bodies.MaxCompressionsPerOneHandshake)+1),
			frame.Frame{
				RequestID: 0,
				Body:      bodies.HandshakeAnswer{Compressions: []fields.Compression{fields.Zstd}},
			},
			false,
		},
		{
			slices.Repeat([]fields.Compression{fields.Zstd}, int(bodies.MaxCompressionsPerOneHandshake)-1),
			frame.Frame{
				RequestID: 0,
				Body:      bodies.HandshakeAnswer{Compressions: []fields.Compression{fields.Zstd}},
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_HandshakeAnswer %v", i),
			func(t *testing.T) {
				f, err := fb.NewHandshakeAnswer(tt.c)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_ErrorAnswer(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		r       fields.RequestID
		e       err.ProtocolError
		want    frame.Frame
		wantErr bool
	}{
		{
			0,
			nil,
			frame.Frame{
				Body: bodies.ErrorAnswer{},
			},
			true,
		},
		{
			1,
			nil,
			frame.Frame{
				RequestID: 1,
				Body:      bodies.ErrorAnswer{},
			},
			true,
		},
		{
			0,
			errs.NewErrorInternalErrorWithTracebackID(nil, fields.TracebackID{}),
			frame.Frame{
				Body: bodies.ErrorAnswer{Err: errs.NewErrorInternalErrorWithTracebackID(nil, fields.TracebackID{})},
			},
			false,
		},
		{
			1,
			errs.NewErrorInternalErrorWithTracebackID(nil, fields.TracebackID{}),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.ErrorAnswer{Err: errs.NewErrorInternalErrorWithTracebackID(nil, fields.TracebackID{})},
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_ErrorAnswer %v", i),
			func(t *testing.T) {
				f, err := fb.NewErrAnswer(tt.r, tt.e)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_Batch(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		r       fields.RequestID
		b       *BatchRequestsBuilder
		want    frame.Frame
		wantErr bool
	}{
		{
			0,
			NewBatchRequestsBuilder(),
			frame.Frame{
				Body: bodies.Batch{},
			},
			true,
		},
		{
			1,
			NewBatchRequestsBuilder(),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Batch{},
			},
			true,
		},
		{
			0,
			NewBatchRequestsBuilder().
				AddRead(fields.Key("qwerty")),
			frame.Frame{
				Body: bodies.Batch{Requests: []bodies.Request{{Number: fields.RequestNumber(0), Body: bodies.Read("qwerty")}}},
			},
			true,
		},
		{
			1,
			NewBatchRequestsBuilder().
				AddRead(fields.Key("")),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Batch{Requests: []bodies.Request{{Number: fields.RequestNumber(0), Body: bodies.Read("")}}},
			},
			true,
		},
		{
			1,
			NewBatchRequestsBuilder().
				AddRead(fields.Key("qwerty")),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.Batch{Requests: []bodies.Request{{Number: fields.RequestNumber(0), Body: bodies.Read("qwerty")}}},
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_ErrorAnswer %v", i),
			func(t *testing.T) {
				f, err := fb.NewBatch(tt.r, *tt.b)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_AnswerBatch(t *testing.T) {
	fb := NewFrameBuilder(1 << 10)

	tests := []struct {
		r       fields.RequestID
		b       *BatchResultsBuilder
		want    frame.Frame
		wantErr bool
	}{
		{
			0,
			NewBatchResultsBuilder(),
			frame.Frame{
				Body: bodies.BatchAnswer{},
			},
			true,
		},
		{
			1,
			NewBatchResultsBuilder(),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.BatchAnswer{},
			},
			true,
		},
		{
			0,
			NewBatchResultsBuilder().
				AddWrite(0),
			frame.Frame{
				Body: bodies.BatchAnswer{{Number: fields.RequestNumber(0), Body: bodies.WriteAnswer{}}},
			},
			true,
		},
		{
			1,
			NewBatchResultsBuilder().
				AddRead(1, values.JSON("qwerty")),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.BatchAnswer{{Number: fields.RequestNumber(1), Body: bodies.ReadAnswer{Value: values.JSON("qwerty")}}},
			},
			true,
		},
		{
			1,
			NewBatchResultsBuilder().
				AddRead(1, values.Int(123)),
			frame.Frame{
				RequestID: 1,
				Body:      bodies.BatchAnswer{{Number: fields.RequestNumber(1), Body: bodies.ReadAnswer{Value: values.Int(123)}}},
			},
			false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestFrameBuilder_BatchAnswer %v", i),
			func(t *testing.T) {
				f, err := fb.NewBatchAnswer(tt.r, *tt.b)
				if err == nil == tt.wantErr {
					t.Fatalf("want err: %v, but got %v", tt.wantErr, err)
					return
				}

				if tt.wantErr {
					return
				}

				if f.RequestID != tt.want.RequestID {
					t.Errorf("request ID: got %v, want %v", f.RequestID, tt.want.RequestID)
				}

				compareBodies(t, f.Body, tt.want.Body)
			},
		)
	}
}

func TestFrameBuilder_LimitExceeded(t *testing.T) {
	fb := NewFrameBuilder(1)

	t.Run(
		"TestFrameBuilder_LimitExceeded",
		func(t *testing.T) {
			_, err := fb.NewPingAnswer(1)
			if err == nil {
				t.Fatal("didn't gtt err")
				return
			}
		},
	)
}
