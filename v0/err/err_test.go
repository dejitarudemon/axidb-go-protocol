package err

import (
	"errors"
	"testing"
)

func TestDecodeError(t *testing.T) {
	sentinel := errors.New("eof")

	tests := []struct {
		name string
		err  DecodeError
		want string
	}{
		{"with source", NewDecodeError("reader eof", sentinel), "reader eof: eof"},
		{"nil source", NewDecodeError("got nil reader", nil), "got nil reader: <nil>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}

	if !errors.Is(NewDecodeError("reader eof", sentinel), sentinel) {
		t.Error("Unwrap() did not return the source error")
	}
}

func TestBuildError(t *testing.T) {
	sentinel := errors.New("too big")

	tests := []struct {
		name string
		err  BuildError
		want string
	}{
		{"with source", NewBuildError("frame too large", sentinel), "frame too large: too big"},
		{"nil source", NewBuildError("frame size is 8, but limit is 1", nil), "frame size is 8, but limit is 1: <nil>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}

	if !errors.Is(NewBuildError("frame too large", sentinel), sentinel) {
		t.Error("Unwrap() did not return the source error")
	}
}

func TestBrokenFrameError(t *testing.T) {
	tests := []struct {
		name string
		got  []byte
		want string
	}{
		{"empty", nil, "got: "},
		{"magic", []byte{0x11, 0xFF}, "got: 11 FF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBrokenFrameError(tt.got).Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	e := NewValidationError("nil Body", "target", "Frame")

	if got, want := e.Error(), "nil Body. use Details() to get more info"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	details := e.Details()
	if len(details) != 2 || details[0] != "target" || details[1] != "Frame" {
		t.Errorf("Details() = %v, want [target Frame]", details)
	}
}

func TestTrailledError(t *testing.T) {
	e := NewTrailledError(8, 4)

	if got, want := e.Error(), "trailled data. expected 8 bytes, but consumed 4"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestMismatchedChecksum(t *testing.T) {
	e := NewMismatchedChecksum(0x6E389900, 0xBD0BE338)

	if got, want := e.Error(), "mismatched checksum: got 6E389900, expected BD0BE338"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
