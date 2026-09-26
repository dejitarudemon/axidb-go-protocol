package bodies

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func TestSimpleOK(t *testing.T) {
	tests := []struct {
		name     string
		s        simpleOK
		wantSize int
		wantErr  bool
	}{
		{"empty", simpleOK{}, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Size(); got != tt.wantSize {
				t.Errorf("Size() = %v, want %v", got, tt.wantSize)
			}

			testutil.AssertErr(t, tt.s.IsValid(), tt.wantErr)
		})
	}
}
