package fields

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func TestType(t *testing.T) {
	tests := []struct {
		typ       Type
		want      []byte
		wantStr   string
		wantValid bool
	}{
		{Bytes, []byte{0x00}, "Bytes", true},
		{TypedArray, []byte{0x01}, "Typed Array", true},
		{UntypedArray, []byte{0x02}, "Untyped Array", true},
		{Int, []byte{0x03}, "Int", true},
		{Uint, []byte{0x04}, "Uint", true},
		{Float, []byte{0x05}, "Float", true},
		{String, []byte{0x06}, "String", true},
		{JSON, []byte{0x07}, "JSON", true},
		{Type(8), []byte{0x08}, "Unknown (8)", false},
		{Type(255), []byte{0xFF}, "Unknown (255)", false},
	}

	for _, tt := range tests {
		t.Run(tt.wantStr, func(t *testing.T) {
			testutil.AssertEncoded(t, tt.typ, tt.want)

			if got := tt.typ.String(); got != tt.wantStr {
				t.Errorf("String() = %q, want %q", got, tt.wantStr)
			}

			if got := tt.typ.IsValid(); got != tt.wantValid {
				t.Errorf("IsValid() = %v, want %v", got, tt.wantValid)
			}
		})
	}
}
