package bodies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func TestBatchReques_Size(t *testing.T) {
	tests := []struct {
		r    Request
		want int
	}{
		{Request{}, 9},
		{Request{0, nil}, 9},
		{Request{1, Read{}}, 9},
		{Request{2, Read("key")}, 12},
		{Request{3, Delete{[]byte("another-key")}}, 20},
		{Request{4, Write{[]byte("key"), nil}}, 16},
		{Request{5, Handshake{"user", [32]byte{}, []fields.Compression{fields.None, fields.Lz4}}}, 52},
		{Request{6, Ping{}}, 9},
		{Request{7, Batch{}}, 14},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatchReques_Size %v", tt.r),
			func(t *testing.T) {
				if got := tt.r.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRequest_Encode(t *testing.T) {
	tests := []struct {
		r    Request
		want []byte
	}{
		{Request{}, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 00}},
		{Request{0, nil}, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{Request{1, Read{}}, []byte{0x00, 0x00, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00}},
		{Request{2, Read("key")}, []byte{0x00, 0x00, 0x00, 0x02, 0x02, 0x00, 0x00, 0x00, 0x03, 0x6B, 0x65, 0x79}},
		{Request{3, Delete{[]byte("another-key")}}, []byte{0x00, 0x00, 0x00, 0x03, 0x04, 0x00, 0x00, 0x00, 0x0B, 0x61, 0x6E, 0x6F, 0x74, 0x68, 0x65, 0x72, 0x2D, 0x6B, 0x65, 0x79}},
		{Request{4, Write{[]byte("key"), nil}}, []byte{0x00, 0x00, 0x00, 0x04, 0x03, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x03, 0x6B, 0x65, 0x79}},
		{Request{5, Handshake{"user", [32]byte{}, []fields.Compression{fields.None, fields.Lz4}}}, []byte{
			0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x2B, 0x00, 0x00, 0x00, 0x04, 0x75, 0x73, 0x65,
			0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x02, 0x00, 0x02,
		},
		},
		{Request{6, Ping{}}, []byte{0x00, 0x00, 0x00, 0x06, 0x06, 0x00, 0x00, 0x00, 0x00}},
		{Request{7, Batch{}}, []byte{0x00, 0x00, 0x00, 0x07, 0x05, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestRequest_Encode %v", tt.r),
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

func TestRequest_IsValid(t *testing.T) {
	tests := []struct {
		r    Request
		want bool
	}{
		{Request{}, true},
		{Request{0, nil}, true},
		{Request{1, Read{}}, true},
		{Request{2, Read("key")}, false},
		{Request{3, Delete{[]byte("another-key")}}, false},
		{Request{4, Write{[]byte("key"), nil}}, true},
		{Request{5, Handshake{"user", [32]byte{}, []fields.Compression{fields.None, fields.Lz4}}}, true},
		{Request{6, Ping{}}, true},
		{Request{7, Batch{}}, true},
		{Request{8, Write{[]byte("key"), values.Int(1)}}, false},
		{Request{9, Write{[]byte("key"), values.JSON([]byte{0x00})}}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestRequest_IsValid %v", tt.r),
			func(t *testing.T) {
				if got := tt.r.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestBatch_Size(t *testing.T) {
	tests := []struct {
		b    Batch
		want int
	}{
		{Batch{}, 5},
		{Batch{false, false, false, []Request{}}, 5},
		{Batch{true, false, false, []Request{}}, 5},
		{Batch{false, true, false, []Request{}}, 5},
		{Batch{false, false, true, []Request{}}, 5},
		{Batch{false, false, true, []Request{{}}}, 14},
		{Batch{false, false, true, []Request{{0, nil}}}, 14},
		{Batch{false, false, true, []Request{{0, nil}}}, 14},
		{Batch{false, false, true, []Request{{0, Read{}}}}, 14},
		{Batch{false, false, true, []Request{{0, Read("key")}}}, 17},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{0, Read("yek")},
		}}, 29},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{1, Read("yek")},
		}}, 29},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{1, Read("yek")},
			{2, Ping{}},
		}}, 38},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatch_Size %v", tt.b),
			func(t *testing.T) {
				if got := tt.b.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestBatch_Encode(t *testing.T) {
	tests := []struct {
		b    Batch
		want []byte
	}{
		{Batch{}, []byte{0x00, 0x00, 0x00, 0x00, 0x00}},
		{Batch{false, false, false, []Request{}}, []byte{0x00, 0x00, 0x00, 0x00, 0x00}},
		{Batch{true, false, false, []Request{}}, []byte{0x01, 0x00, 0x00, 0x00, 0x00}},
		{Batch{false, true, false, []Request{}}, []byte{0x02, 0x00, 0x00, 0x00, 0x00}},
		{Batch{false, false, true, []Request{}}, []byte{0x04, 0x00, 0x00, 0x00, 0x00}},
		{Batch{false, false, true, []Request{{}}}, []byte{0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{Batch{false, false, true, []Request{{0, nil}}}, []byte{0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{Batch{false, false, true, []Request{{0, nil}}}, []byte{0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{Batch{false, false, true, []Request{{0, Read{}}}}, []byte{0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00}},
		{Batch{false, false, true, []Request{{0, Read("key")}}}, []byte{0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x6B, 0x65, 0x79}},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{0, Read("yek")},
		}}, []byte{0x04, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x6B, 0x65, 0x79, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x79, 0x65, 0x6B}},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{1, Read("yek")},
		}}, []byte{0x04, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x6B, 0x65, 0x79, 0x00, 0x00, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x03, 0x79, 0x65, 0x6B}},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{1, Read("yek")},
			{2, Ping{}},
		}},
			[]byte{
				0x04, 0x00, 0x00, 0x00, 0x03,
				0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x6B, 0x65, 0x79,
				0x00, 0x00, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x03, 0x79, 0x65, 0x6B,
				0x00, 0x00, 0x00, 0x02, 0x06, 0x00, 0x00, 0x00, 0x00,
			}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatch_Encode %v", tt.b),
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

func TestBatch_Command(t *testing.T) {
	tests := []struct {
		b    Batch
		want fields.Command
	}{
		{Batch{}, fields.Batch},
		{Batch{false, false, false, []Request{}}, fields.Batch},
		{Batch{true, false, false, []Request{}}, fields.Batch},
		{Batch{false, true, false, []Request{}}, fields.Batch},
		{Batch{false, false, true, []Request{}}, fields.Batch},
		{Batch{false, false, true, []Request{{}}}, fields.Batch},
		{Batch{false, false, true, []Request{{0, nil}}}, fields.Batch},
		{Batch{false, false, true, []Request{{0, nil}}}, fields.Batch},
		{Batch{false, false, true, []Request{{0, Read{}}}}, fields.Batch},
		{Batch{false, false, true, []Request{{0, Read("key")}}}, fields.Batch},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{0, Read("yek")},
		}}, fields.Batch},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{1, Read("yek")},
		}}, fields.Batch},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{1, Read("yek")},
			{2, Ping{}},
		}}, fields.Batch},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatch_Command %v", tt.b),
			func(t *testing.T) {
				if got := tt.b.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestBatch_IsValid(t *testing.T) {
	tests := []struct {
		b    Batch
		want bool
	}{
		{Batch{}, true},
		{Batch{false, false, false, []Request{}}, true},
		{Batch{true, false, false, []Request{}}, true},
		{Batch{false, true, false, []Request{}}, true},
		{Batch{false, false, true, []Request{}}, true},
		{Batch{false, false, true, []Request{{}}}, true},
		{Batch{false, false, true, []Request{{0, nil}}}, true},
		{Batch{false, false, true, []Request{{0, Read{}}}}, true},
		{Batch{false, false, true, []Request{{0, Read("key")}}}, false},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{0, Read("yek")},
		}}, true},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{1, Read("yek")},
		}}, false},
		{Batch{false, false, true, []Request{
			{0, Read("key")},
			{1, Read("yek")},
			{2, Ping{}},
		}}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestBatch_IsValid %v", tt.b),
			func(t *testing.T) {
				if got := tt.b.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
