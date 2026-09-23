package fields

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestError_Size(t *testing.T) {
	tests := []struct {
		e    Error
		want int
	}{
		{NoHello, 2},
		{UnsupportedVersion, 2},
		{UnexpectedCommand, 2},
		{UnsupportedCommand, 2},
		{RequestsConflict, 2},
		{UnsupportedCompression, 2},
		{BodyLimitIsExceeded, 2},
		{MismatchedChecksum, 2},
		{InternalError, 2},
		{MalformedValue, 2},
		{NotFound, 2},
		{ProhibitedCompression, 2},
		{RequestInterrupted, 2},
		{BatchLimitIsExceeded, 2},
		{UnexpectedCommandInBatch, 2},
		{InvalidRequestID, 2},
		{Unauthorized, 2},
		{RestrictedRequest, 2},
		{Error(19), 2},
		{Error(255), 2},
		{Error(65535), 2},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.e),
			func(t *testing.T) {
				if got := tt.e.Size(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestError_Encode(t *testing.T) {
	tests := []struct {
		e    Error
		want []byte
	}{
		{NoHello, []byte{0x00, 0x00}},
		{UnsupportedVersion, []byte{0x00, 0x01}},
		{UnexpectedCommand, []byte{0x00, 0x02}},
		{UnsupportedCommand, []byte{0x00, 0x03}},
		{RequestsConflict, []byte{0x00, 0x04}},
		{UnsupportedCompression, []byte{0x00, 0x05}},
		{BodyLimitIsExceeded, []byte{0x00, 0x06}},
		{MismatchedChecksum, []byte{0x00, 0x07}},
		{InternalError, []byte{0x00, 0x08}},
		{MalformedValue, []byte{0x00, 0x09}},
		{NotFound, []byte{0x00, 0x0A}},
		{ProhibitedCompression, []byte{0x00, 0x0B}},
		{RequestInterrupted, []byte{0x00, 0x0C}},
		{BatchLimitIsExceeded, []byte{0x00, 0x0D}},
		{UnexpectedCommandInBatch, []byte{0x00, 0x0E}},
		{InvalidRequestID, []byte{0x00, 0x0F}},
		{Unauthorized, []byte{0x00, 0x10}},
		{RestrictedRequest, []byte{0x00, 0x11}},
		{Error(19), []byte{0x00, 0x13}},
		{Error(255), []byte{0x00, 0xFF}},
		{Error(65535), []byte{0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.e),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.e.Size())

				tt.e.Encode(&buf)
				got := buf.Bytes()

				if len(got) != tt.e.Size() {
					t.Fatalf("expected %v bytes, got %v bytes", tt.e.Size(), len(got))
				}

				if !bytes.Equal(got, tt.want) {
					t.Errorf("Encode() = %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestError_String(t *testing.T) {
	tests := []struct {
		e    Error
		want string
	}{
		{NoHello, "No Hello"},
		{UnexpectedCommand, "Unexpected Command"},
		{UnsupportedCommand, "Unsupported Command"},
		{RequestsConflict, "Requests Conflict"},
		{UnsupportedCompression, "Unsupported Compression"},
		{BodyLimitIsExceeded, "Body Limit Is Exceeded"},
		{MismatchedChecksum, "Mismatched Checksum"},
		{InternalError, "Internal Error"},
		{MalformedValue, "Malformed Value"},
		{NotFound, "Not Found"},
		{ProhibitedCompression, "Prohibited Compression"},
		{RequestInterrupted, "Request Interrupted"},
		{BatchLimitIsExceeded, "Batch Limit Is Exceeded"},
		{UnexpectedCommandInBatch, "Unexpected Command In Batch"},
		{InvalidRequestID, "Invalid Request ID"},
		{Unauthorized, "Unauthorized"},
		{RestrictedRequest, "Restricted Request"},
		{Error(19), "Unknown (19)"},
		{Error(255), "Unknown (255)"},
		{Error(65535), "Unknown (65535)"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.e),
			func(t *testing.T) {
				if got := tt.e.String(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestError_IsValid(t *testing.T) {
	tests := []struct {
		e    Error
		want bool
	}{
		{NoHello, true},
		{UnexpectedCommand, true},
		{UnsupportedCommand, true},
		{RequestsConflict, true},
		{UnsupportedCompression, true},
		{BodyLimitIsExceeded, true},
		{MismatchedChecksum, true},
		{InternalError, true},
		{MalformedValue, true},
		{NotFound, true},
		{ProhibitedCompression, true},
		{RequestInterrupted, true},
		{BatchLimitIsExceeded, true},
		{UnexpectedCommandInBatch, true},
		{InvalidRequestID, true},
		{Unauthorized, true},
		{RestrictedRequest, true},
		{Error(19), false},
		{Error(255), false},
		{Error(65535), false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.e),
			func(t *testing.T) {
				if got := tt.e.IsValid(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
