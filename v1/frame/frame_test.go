package frame

import (
	"fmt"
	"testing"
)

func TestSize_Frame(t *testing.T) {
	tests := []struct {
		f    Frame
		want int
	}{}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSize_Frame %v", tt.f),
			func(t *testing.T) {
				if got := tt.f.Size(); got != tt.want {
				}
			},
		)
	}
}
