package bodies

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func TestBatchAnswerReques_Size(t *testing.T) {
	tests := []struct {
		r    Result
		want int
	}{
		{Result{}, 8},
		{Result{0, nil}, 8},
		{Result{1, Read{}}, 8},
		{Result{2, Read{[]byte("key")}}, 11},
		{Result{3, Write{[]byte("key"), nil}}, 15},
		{Result{4, Ping{}}, 8},
		{Result{5, Batch{}}, 13},
		{Result{6, WriteAnswer{}}, 9},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatchAnswerReques_Size %v", tt.r),
			func(t *testing.T) {
				if got := tt.r.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestResult_Encode(t *testing.T) {
	tests := []struct {
		r    Result
		want []byte
	}{
		{
			Result{},
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			Result{0, nil},
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			Result{1, Read{}},
			[]byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
		},
		{
			Result{2, Read{[]byte("key")}},
			[]byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x6B, 0x65, 0x79},
		},
		{
			Result{3, Write{[]byte("key"), nil}},
			[]byte{0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x03, 0x6B, 0x65, 0x79},
		},
		{
			Result{4, Ping{}},
			[]byte{0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00},
		},
		{
			Result{5, Batch{}},
			[]byte{0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			Result{6, WriteAnswer{}},
			[]byte{0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x01, 0x01},
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestResult_Encode %v", tt.r),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.r.Size())

				tt.r.Encode(&buf)

				got := buf.Bytes()

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestResult_IsValid(t *testing.T) {
	tests := []struct {
		r    Result
		want bool
	}{
		{Result{}, true},
		{Result{0, nil}, true},
		{Result{1, Read{}}, true},
		{Result{2, Read{[]byte("key")}}, true},
		{Result{3, Write{[]byte("key"), nil}}, true},
		{Result{4, Ping{}}, true},
		{Result{5, Batch{}}, true},
		{Result{6, WriteAnswer{}}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestResult_IsValid %v", tt.r),
			func(t *testing.T) {
				if got := tt.r.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestBatchAnswer_Size(t *testing.T) {
	tests := []struct {
		b    BatchAnswer
		want int
	}{
		{BatchAnswer{}, 5},
		{BatchAnswer{[]Result{}}, 5},
		{BatchAnswer{[]Result{{}}}, 13},
		{BatchAnswer{[]Result{{0, nil}}}, 13},
		{BatchAnswer{[]Result{{0, Read{}}}}, 13},
		{BatchAnswer{[]Result{{0, WriteAnswer{}}}}, 14},
		{BatchAnswer{[]Result{
			{0, WriteAnswer{}},
			{0, DeleteAnswer{}},
		}}, 23},
		{BatchAnswer{[]Result{
			{0, WriteAnswer{}},
			{1, DeleteAnswer{}},
		}}, 23},
		{BatchAnswer{[]Result{
			{1, ReadAnswer{Value: values.String("message")}},
			{2, WriteAnswer{}},
		}}, 35},
		{BatchAnswer{[]Result{
			{1, ErrorAnswer{errs.NewErrorInternalErrorWithTracebackID(errors.New(""), [16]byte{0xDA, 0xE1, 0x2F, 0x25, 0x35, 0x1B, 0x48, 0x75, 0x99, 0xDE, 0x1D, 0xF5, 0x2A, 0x76, 0x3F, 0x96})}},
			{2, ErrorAnswer{errs.NewErrorRequestInterruptedWithTracebackID(2, [16]byte{0x77, 0xA5, 0x7D, 0x8E, 0xCC, 0x20, 0x40, 0x7F, 0x8C, 0x65, 0x5E, 0x94, 0x5F, 0xE2, 0xA8, 0x91})}},
		}}, 59},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatchAnswer_Size %v", tt.b),
			func(t *testing.T) {
				if got := tt.b.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestBatchAnswer_Encode(t *testing.T) {
	tests := []struct {
		b    BatchAnswer
		want []byte
	}{
		{BatchAnswer{}, []byte{0x01, 0x00, 0x00, 0x00, 0x00}},
		{BatchAnswer{[]Result{}}, []byte{0x01, 0x00, 0x00, 0x00, 0x00}},
		{
			BatchAnswer{[]Result{{}}},
			[]byte{0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			BatchAnswer{[]Result{{0, nil}}},
			[]byte{0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			BatchAnswer{[]Result{{0, Read{}}}},
			[]byte{0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			BatchAnswer{[]Result{{0, WriteAnswer{}}}},
			[]byte{0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01},
		},
		{
			BatchAnswer{[]Result{
				{0, WriteAnswer{}},
				{0, DeleteAnswer{}},
			}},
			[]byte{0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01},
		},
		{
			BatchAnswer{[]Result{
				{0, WriteAnswer{}},
				{1, DeleteAnswer{}},
			}},
			[]byte{0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01},
		},
		{
			BatchAnswer{[]Result{
				{1, ReadAnswer{Value: values.String("message")}},
				{2, WriteAnswer{}},
			}},
			[]byte{
				0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00,
				0x00, 0x00, 0x0D, 0x01, 0x06, 0x00, 0x00, 0x00, 0x07, 0x6D,
				0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x00, 0x00, 0x00, 0x02,
				0x00, 0x00, 0x00, 0x01, 0x01},
		},
		{
			BatchAnswer{[]Result{
				{1, ErrorAnswer{errs.NewErrorInternalErrorWithTracebackID(errors.New(""), [16]byte{0xDA, 0xE1, 0x2F, 0x25, 0x35, 0x1B, 0x48, 0x75, 0x99, 0xDE, 0x1D, 0xF5, 0x2A, 0x76, 0x3F, 0x96})}},
				{2, ErrorAnswer{errs.NewErrorRequestInterruptedWithTracebackID(2, [16]byte{0x77, 0xA5, 0x7D, 0x8E, 0xCC, 0x20, 0x40, 0x7F, 0x8C, 0x65, 0x5E, 0x94, 0x5F, 0xE2, 0xA8, 0x91})}},
			}},
			[]byte{
				0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00,
				0x00, 0x00, 0x13, 0x00, 0x00, 0x08, 0xDA, 0xE1, 0x2F, 0x25,
				0x35, 0x1B, 0x48, 0x75, 0x99, 0xDE, 0x1D, 0xF5, 0x2A, 0x76,
				0x3F, 0x96, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x13,
				0x00, 0x00, 0x0C, 0x77, 0xA5, 0x7D, 0x8E, 0xCC, 0x20, 0x40,
				0x7F, 0x8C, 0x65, 0x5E, 0x94, 0x5F, 0xE2, 0xA8, 0x91,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatchAnswer_Encode %v", tt.b),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.b.Size())

				tt.b.Encode(&buf)

				got := buf.Bytes()

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestBatchAnswer_Command(t *testing.T) {
	tests := []struct {
		b    BatchAnswer
		want fields.Command
	}{
		{BatchAnswer{}, fields.Answer},
		{BatchAnswer{[]Result{}}, fields.Answer},
		{BatchAnswer{[]Result{{}}}, fields.Answer},
		{BatchAnswer{[]Result{{0, nil}}}, fields.Answer},
		{BatchAnswer{[]Result{{0, Read{}}}}, fields.Answer},
		{BatchAnswer{[]Result{{0, WriteAnswer{}}}}, fields.Answer},
		{BatchAnswer{[]Result{
			{0, WriteAnswer{}},
			{0, DeleteAnswer{}},
		}}, fields.Answer},
		{BatchAnswer{[]Result{
			{0, WriteAnswer{}},
			{1, DeleteAnswer{}},
		}}, fields.Answer},
		{BatchAnswer{[]Result{
			{1, ReadAnswer{Value: values.String("message")}},
			{2, WriteAnswer{}},
		}}, fields.Answer},
		{BatchAnswer{[]Result{
			{1, ErrorAnswer{errs.NewErrorInternalErrorWithTracebackID(errors.New(""), [16]byte{0xDA, 0xE1, 0x2F, 0x25, 0x35, 0x1B, 0x48, 0x75, 0x99, 0xDE, 0x1D, 0xF5, 0x2A, 0x76, 0x3F, 0x96})}},
			{2, ErrorAnswer{errs.NewErrorRequestInterruptedWithTracebackID(2, [16]byte{0x77, 0xA5, 0x7D, 0x8E, 0xCC, 0x20, 0x40, 0x7F, 0x8C, 0x65, 0x5E, 0x94, 0x5F, 0xE2, 0xA8, 0x91})}},
		}}, fields.Answer},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatchAnswer_Command %v", tt.b),
			func(t *testing.T) {
				if got := tt.b.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestBatchAnswer_IsValid(t *testing.T) {
	tests := []struct {
		b    BatchAnswer
		want bool
	}{
		{BatchAnswer{}, true},
		{BatchAnswer{[]Result{}}, true},
		{BatchAnswer{[]Result{{}}}, true},
		{BatchAnswer{[]Result{{0, nil}}}, true},
		{BatchAnswer{[]Result{{0, Read{}}}}, true},
		{BatchAnswer{[]Result{{0, WriteAnswer{}}}}, false},
		{BatchAnswer{[]Result{
			{0, WriteAnswer{}},
			{0, DeleteAnswer{}},
		}}, true},
		{BatchAnswer{[]Result{
			{0, WriteAnswer{}},
			{1, DeleteAnswer{}},
		}}, false},
		{BatchAnswer{[]Result{
			{1, ReadAnswer{Value: values.String("message")}},
			{2, WriteAnswer{}},
		}}, false},
		{BatchAnswer{[]Result{
			{1, ErrorAnswer{errs.NewErrorInternalErrorWithTracebackID(errors.New(""), [16]byte{0xDA, 0xE1, 0x2F, 0x25, 0x35, 0x1B, 0x48, 0x75, 0x99, 0xDE, 0x1D, 0xF5, 0x2A, 0x76, 0x3F, 0x96})}},
			{2, ErrorAnswer{errs.NewErrorRequestInterruptedWithTracebackID(2, [16]byte{0x77, 0xA5, 0x7D, 0x8E, 0xCC, 0x20, 0x40, 0x7F, 0x8C, 0x65, 0x5E, 0x94, 0x5F, 0xE2, 0xA8, 0x91})}},
		}}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatchAnswer_IsValid %v", tt.b),
			func(t *testing.T) {
				if got := tt.b.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
