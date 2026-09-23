package bodies

import (
	"fmt"
	"testing"
)

func TestSimpleOK_Size(t *testing.T) {
	tests := []struct {
		s    simpleOK
		want int
	}{
		{simpleOK{}, 2},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSimpleOK_Size %v", tt.s),
			func(t *testing.T) {
				if got := tt.s.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestSimpleOK_IsValid(t *testing.T) {
	tests := []struct {
		s    simpleOK
		want bool
	}{
		{simpleOK{}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSimpleOK_IsValid %v", tt.s),
			func(t *testing.T) {
				if got := tt.s.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
