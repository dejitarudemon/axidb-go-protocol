package values

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

var (
	json1, _ = json.Marshal(map[string]string{})
	json2, _ = json.Marshal(map[string]string{"1": "2"})
	json3, _ = json.Marshal(map[string]any{"1": 1})
)

func TestJSON_Size(t *testing.T) {
	tests := []struct {
		c    JSON
		want int
	}{
		{JSON(json1), 6},
		{JSON(json2), 13},
		{JSON(json3), 11},
		{JSON([]byte{}), 4},
		{JSON([]byte{0x00, 0x01, 0x02}), 7},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.Size(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestJSON_Encode(t *testing.T) {
	tests := []struct {
		c    JSON
		want []byte
	}{
		{JSON(json1), []byte{0x00, 0x00, 0x00, 0x02, 0x7B, 0x7D}},
		{JSON(json2), []byte{0x00, 0x00, 0x00, 0x09, 0x7B, 0x22, 0x31, 0x22, 0x3A, 0x22, 0x32, 0x22, 0x7D}},
		{JSON(json3), []byte{0x00, 0x00, 0x00, 0x07, 0x7B, 0x22, 0x31, 0x22, 0x3A, 0x31, 0x7D}},
		{JSON([]byte{}), []byte{0x00, 0x00, 0x00, 0x00}},
		{JSON([]byte{0x00, 0x01, 0x02}), []byte{0x00, 0x00, 0x00, 0x03, 0x00, 0x01, 0x02}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.c.Size())

				tt.c.Encode(&buf)
				got := buf.Bytes()

				if len(got) != tt.c.Size() {
					t.Fatalf("expected %v bytes, got %v bytes", tt.c.Size(), len(got))
				}

				if !bytes.Equal(got, tt.want) {
					t.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestJSON_Type(t *testing.T) {
	tests := []struct {
		c    JSON
		want fields.Type
	}{
		{JSON(json1), fields.JSON},
		{JSON(json2), fields.JSON},
		{JSON(json3), fields.JSON},
		{JSON([]byte{}), fields.JSON},
		{JSON([]byte{0x00, 0x01, 0x02}), fields.JSON},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				got := tt.c.Type()

				if got != tt.c.Type() {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestJSON_IsValid(t *testing.T) {
	tests := []struct {
		c       JSON
		wantErr bool
	}{
		{JSON(json1), false},
		{JSON(json2), false},
		{JSON(json3), false},
		{JSON([]byte{}), true},
		{JSON([]byte{0x00, 0x01, 0x02}), true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.IsValid(); got == nil == tt.wantErr {
					t.Errorf("got %v, want %v", got, tt.wantErr)
				}
			},
		)
	}
}
