package bodies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func TestReadAnswer_Size(t *testing.T) {
	tests := []struct {
		r    ReadAnswer
		want int
	}{
		{ReadAnswer{}, 2},
		{ReadAnswer{values.Int(0)}, 11},
		{ReadAnswer{values.Bytes("data")}, 11},
		{ReadAnswer{values.String("some-data")}, 16},
		{
			ReadAnswer{
				values.UntypedArray{

					values.Int(1),
					values.String("hello world"),
					values.Float(2.1),
				},
			}, 41,
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestReadAnswer_Size %v", tt.r),
			func(t *testing.T) {
				if got := tt.r.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestReadAnswer_Encode(t *testing.T) {
	tests := []struct {
		r    ReadAnswer
		want []byte
	}{
		{ReadAnswer{}, []byte{0x01, 0x02}},
		{ReadAnswer{values.Int(0)}, []byte{0x01, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{ReadAnswer{values.Bytes("data")}, []byte{0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x04, 0x64, 0x61, 0x74, 0x61}},
		{ReadAnswer{values.String("some-data")}, []byte{0x01, 0x02, 0x06, 0x00, 0x00, 0x00, 0x09, 0x73, 0x6F, 0x6D, 0x65, 0x2D, 0x64, 0x61, 0x74, 0x61}},
		{
			ReadAnswer{
				values.UntypedArray{

					values.Int(1),
					values.String("hello world"),
					values.Float(2.1),
				},
			},
			[]byte{
				0x01, 0x02, 0x02, 0x00, 0x00, 0x00, 0x03, 0x03, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x06, 0x00, 0x00, 0x00,
				0x0B, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x77, 0x6F, 0x72,
				0x6C, 0x64, 0x05, 0x40, 0x00, 0xCC, 0xCC, 0xCC, 0xCC, 0xCC,
				0xCD,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestReadAnswer_Encode %v", tt.r),
			func(t *testing.T) {
				buf := buffer.Slice{}
				buf.Preallocate(tt.r.Size())

				tt.r.Encode(&buf)

				got := buf.Bytes()

				if len(got) != len(tt.want) {
					t.Fatalf("got %v want %v", len(got), len(tt.want))
				}

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestReadAnswer_Command(t *testing.T) {
	tests := []struct {
		r    ReadAnswer
		want fields.Command
	}{
		{ReadAnswer{}, fields.Answer},
		{ReadAnswer{values.Int(0)}, fields.Answer},
		{ReadAnswer{values.Bytes("data")}, fields.Answer},
		{ReadAnswer{values.String("some-data")}, fields.Answer},
		{
			ReadAnswer{
				values.UntypedArray{

					values.Int(1),
					values.String("hello world"),
					values.Float(2.1),
				},
			}, fields.Answer,
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestReadAnswer_Command %v", tt.r),
			func(t *testing.T) {
				if got := tt.r.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestReadAnswer_IsValid(t *testing.T) {
	tests := []struct {
		r    ReadAnswer
		want bool
	}{
		{ReadAnswer{}, true},
		{ReadAnswer{values.Int(0)}, false},
		{ReadAnswer{values.Bytes("data")}, false},
		{ReadAnswer{values.String("some-data")}, false},
		{
			ReadAnswer{
				values.UntypedArray{

					values.Int(1),
					values.String("hello world"),
					values.Float(2.1),
				},
			}, false,
		},
	}
	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestReadAnswer_IsValid %v", tt.r),
			func(t *testing.T) {
				if got := tt.r.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestReadAnswer_IsResponseTo(t *testing.T) {
	tests := []struct {
		r    ReadAnswer
		want fields.Command
	}{
		{ReadAnswer{}, fields.Read},
		{ReadAnswer{values.Int(0)}, fields.Read},
		{ReadAnswer{values.Bytes("data")}, fields.Read},
		{ReadAnswer{values.String("some-data")}, fields.Read},
		{
			ReadAnswer{
				values.UntypedArray{

					values.Int(1),
					values.String("hello world"),
					values.Float(2.1),
				},
			}, fields.Read,
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestRead_IsResponseTo %v", tt.r),
			func(t *testing.T) {
				if got := tt.r.IsResponseTo(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
