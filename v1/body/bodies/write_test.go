package bodies

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func TestWrite_Size(t *testing.T) {
	tests := []struct {
		w    Write
		want int
	}{
		{Write{}, 4},
		{Write{Key: []byte{0x00}}, 5},
		{Write{Value: nil}, 4},
		{Write{[]byte{0x00, 0x01}, values.Int(1)}, 15},
		{Write{[]byte{0x00}, values.Bytes{}}, 10},
		{Write{[]byte{0x0A, 0x0B}, values.Bytes([]byte{0x01, 0x01, 0x02, 0x03})}, 15},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestWrite_Size %v", tt.w),
			func(t *testing.T) {
				if got := tt.w.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestWrite_Encode(t *testing.T) {
	tests := []struct {
		w    Write
		want []byte
	}{
		{
			Write{},
			[]byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			Write{Key: []byte{0x00}},
			[]byte{0x00, 0x00, 0x00, 0x01, 0x00},
		},
		{
			Write{Value: nil},
			[]byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			Write{[]byte{0x00, 0x01}, values.Int(1)},
			[]byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x01, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			Write{[]byte{0x00}, values.Bytes{}},
			[]byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			Write{[]byte{0x0A, 0x0B}, values.Bytes([]byte{0x01, 0x01, 0x02, 0x03})},
			[]byte{0x00, 0x00, 0x00, 0x02, 0x0A, 0x0B, 0x00, 0x00, 0x00, 0x00, 0x04, 0x01, 0x01, 0x02, 0x03},
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestWrite_Encode %v", tt.w),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.w.Size())

				tt.w.Encode(&buf)

				got := buf.Bytes()

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestWrite_Command(t *testing.T) {
	tests := []struct {
		w    Write
		want command.Code
	}{
		{Write{}, command.Write},
		{Write{Key: []byte{0x00}}, command.Write},
		{Write{Value: nil}, command.Write},
		{Write{[]byte{0x00, 0x01}, values.Int(1)}, command.Write},
		{Write{[]byte{0x00}, values.Bytes{}}, command.Write},
		{Write{[]byte{0x0A, 0x0B}, values.Bytes([]byte{0x01, 0x01, 0x02, 0x03})}, command.Write},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestWrite_Command %v", tt.w),
			func(t *testing.T) {
				if got := tt.w.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestWrite_IsValid(t *testing.T) {
	tests := []struct {
		w    Write
		want bool
	}{
		{Write{}, true},
		{Write{Key: []byte{0x00}}, true},
		{Write{Value: nil}, true},
		{Write{[]byte{0x00, 0x01}, values.Int(1)}, false},
		{Write{[]byte{0x00}, values.Bytes{}}, false},
		{Write{[]byte{0x0A, 0x0B}, values.Bytes([]byte{0x01, 0x01, 0x02, 0x03})}, false},
		{Write{[]byte{0x0A, 0x0B}, values.JSON([]byte{0x01, 0x01, 0x02, 0x03})}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestWrite_IsValid %v", tt.w),
			func(t *testing.T) {
				if got := tt.w.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
