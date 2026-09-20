package err

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func TestCode_Size(t *testing.T) {
	tests := []struct {
		c    Code
		want int
	}{
		{NoHello, 2},
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
		{UnknownErrorCode, 2},
		{Code(18), 2},
		{Code(255), 2},
		{Code(65535), 2},
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

func TestCode_Encode(t *testing.T) {
	tests := []struct {
		c    Code
		want []byte
	}{
		{NoHello, []byte{0x00, 0x00}},
		{UnexpectedCommand, []byte{0x00, 0x01}},
		{UnsupportedCommand, []byte{0x00, 0x02}},
		{RequestsConflict, []byte{0x00, 0x03}},
		{UnsupportedCompression, []byte{0x00, 0x04}},
		{BodyLimitIsExceeded, []byte{0x00, 0x05}},
		{MismatchedChecksum, []byte{0x00, 0x06}},
		{InternalError, []byte{0x00, 0x07}},
		{MalformedValue, []byte{0x00, 0x08}},
		{NotFound, []byte{0x00, 0x09}},
		{ProhibitedCompression, []byte{0x00, 0x0A}},
		{RequestInterrupted, []byte{0x00, 0x0B}},
		{BatchLimitIsExceeded, []byte{0x00, 0x0C}},
		{UnexpectedCommandInBatch, []byte{0x00, 0x0D}},
		{InvalidRequestID, []byte{0x00, 0x0E}},
		{Unauthorized, []byte{0x00, 0x0F}},
		{RestrictedRequest, []byte{0x00, 0x10}},
		{UnknownErrorCode, []byte{0x00, 0x11}},
		{Code(18), []byte{0x00, 0x12}},
		{Code(255), []byte{0x00, 0xFF}},
		{Code(65535), []byte{0xFF, 0xFF}},
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

func TestCode_String(t *testing.T) {
	tests := []struct {
		c    Code
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
		{UnknownErrorCode, "Unknown Error Code"},
		{Code(18), "Unknown (18)"},
		{Code(255), "Unknown (255)"},
		{Code(65535), "Unknown (65535)"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.String(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestCode_IsValid(t *testing.T) {
	tests := []struct {
		c    Code
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
		{UnknownErrorCode, true},
		{Code(18), false},
		{Code(255), false},
		{Code(65535), false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt.c),
			func(t *testing.T) {
				if got := tt.c.IsValid(); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
