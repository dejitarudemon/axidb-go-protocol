package fields

import (
	"math"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func TestError(t *testing.T) {
	tests := []struct {
		e         Error
		want      []byte
		wantStr   string
		wantValid bool
	}{
		{NoHello, []byte{0x00, 0x00}, "No Hello", true},
		{UnsupportedVersion, []byte{0x00, 0x01}, "Unsupported Version", true},
		{UnexpectedCommand, []byte{0x00, 0x02}, "Unexpected Command", true},
		{UnsupportedCommand, []byte{0x00, 0x03}, "Unsupported Command", true},
		{RequestsConflict, []byte{0x00, 0x04}, "Requests Conflict", true},
		{UnsupportedCompression, []byte{0x00, 0x05}, "Unsupported Compression", true},
		{BodyLimitIsExceeded, []byte{0x00, 0x06}, "Body Limit Is Exceeded", true},
		{MismatchedChecksum, []byte{0x00, 0x07}, "Mismatched Checksum", true},
		{InternalError, []byte{0x00, 0x08}, "Internal Error", true},
		{MalformedValue, []byte{0x00, 0x09}, "Malformed Value", true},
		{NotFound, []byte{0x00, 0x0A}, "Not Found", true},
		{ProhibitedCompression, []byte{0x00, 0x0B}, "Prohibited Compression", true},
		{RequestInterrupted, []byte{0x00, 0x0C}, "Request Interrupted", true},
		{BatchLimitIsExceeded, []byte{0x00, 0x0D}, "Batch Limit Is Exceeded", true},
		{UnexpectedCommandInBatch, []byte{0x00, 0x0E}, "Unexpected Command In Batch", true},
		{InvalidRequestID, []byte{0x00, 0x0F}, "Invalid Request ID", true},
		{Unauthorized, []byte{0x00, 0x10}, "Unauthorized", true},
		{RestrictedRequest, []byte{0x00, 0x11}, "Restricted Request", true},
		{Error(18), []byte{0x00, 0x12}, "Unknown (18)", false},
		{Error(255), []byte{0x00, 0xFF}, "Unknown (255)", false},
		{Error(math.MaxUint16), []byte{0xFF, 0xFF}, "Unknown (65535)", false},
	}

	for _, tt := range tests {
		t.Run(tt.wantStr, func(t *testing.T) {
			testutil.AssertEncoded(t, tt.e, tt.want)

			if got := tt.e.String(); got != tt.wantStr {
				t.Errorf("String() = %q, want %q", got, tt.wantStr)
			}

			if got := tt.e.IsValid(); got != tt.wantValid {
				t.Errorf("IsValid() = %v, want %v", got, tt.wantValid)
			}
		})
	}
}
