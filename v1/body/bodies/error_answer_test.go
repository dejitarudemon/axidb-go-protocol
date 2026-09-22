package bodies

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func TestErrorAnswer_Size(t *testing.T) {
	tests := []struct {
		e    ErrorAnswer
		want int
	}{
		{ErrorAnswer{}, 1},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(1, 0)}, 23},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(0, 1)}, 23},
		{ErrorAnswer{errs.NewErrorBodyLimitIsExceeded(1, 0)}, 23},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(0, 1)}, 23},
		{ErrorAnswer{errs.NewErrorInternalError(errors.New(""))}, 19},
		{ErrorAnswer{errs.NewErrorInternalError(nil)}, 19},
		{ErrorAnswer{errs.NewErrorInvalidRequestID(fields.Handshake, 1)}, 19},
		{ErrorAnswer{errs.NewErrorInvalidRequestID(fields.Handshake, 0)}, 19},
		{ErrorAnswer{errs.NewErrorMalformedValue("some-error")}, 29},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(0, 1)}, 19},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(1, 1)}, 19},
		{ErrorAnswer{errs.NewErrorNoHello()}, 19},
		{ErrorAnswer{errs.NewErrorNotFound([]byte("key"))}, 19},
		{ErrorAnswer{errs.NewErrorNotFound([]byte(""))}, 19},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(fields.S2, fields.Handshake)}, 19},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(fields.None, fields.Handshake)}, 19},
		{ErrorAnswer{errs.NewErrorRequestsConflict(0)}, 19},
		{ErrorAnswer{errs.NewErrorRequestInterrupted(0)}, 19},
		{ErrorAnswer{errs.NewErrorRestrictedRequest([]byte("user"), []byte("1.1.1.1"), []byte("key"), fields.Read, 0)}, 19},
		{ErrorAnswer{errs.NewErrorUnauthorized([]byte("1.1.1.1"))}, 19},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(fields.Answer, 0)}, 23},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(fields.Read, 0)}, 23},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(fields.Handshake, fields.Read)}, 20},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(fields.Handshake, fields.Handshake)}, 20},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(fields.Command(255))}, 19},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(fields.Handshake)}, 19},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(fields.S2)}, 19},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(fields.None)}, 19},
		{ErrorAnswer{errs.NewErrorUnsupportedVersion(3)}, 19},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestErrorAnswer_Size %v", tt.e),
			func(t *testing.T) {
				if got := tt.e.Size(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestErrorAnswer_Encode(t *testing.T) {
	tests := []struct {
		e    ErrorAnswer
		want []byte
	}{
		{ErrorAnswer{}, []byte{0x00}},
		{
			ErrorAnswer{errs.NewErrorBatchLimitIsExceededWithTracebackID(1, 0, [16]byte{0x00})},
			[]byte{0x00, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorBatchLimitIsExceededWithTracebackID(0, 1, [16]byte{0x01})},
			[]byte{0x00, 0x00, 0x0D, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			ErrorAnswer{errs.NewErrorBodyLimitIsExceededWithTracebackID(1, 0, [16]byte{0x01, 0x02})},
			[]byte{0x00, 0x00, 0x06, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorBodyLimitIsExceededWithTracebackID(0, 1, [16]byte{0x01, 0x02, 0x03})},
			[]byte{0x00, 0x00, 0x06, 0x01, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			ErrorAnswer{errs.NewErrorInternalErrorWithTracebackID(errors.New(""), [16]byte{0x01, 0x02})},
			[]byte{0x00, 0x00, 0x08, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorInternalErrorWithTracebackID(nil, [16]byte{0x01})},
			[]byte{0x00, 0x00, 0x08, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorInvalidRequestIDWithTracebackID(fields.Handshake, 1, [16]byte{0x01})},
			[]byte{0x00, 0x00, 0x0F, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorInvalidRequestIDWithTracebackID(fields.Handshake, 0, [16]byte{0x02, 0x03})},
			[]byte{0x00, 0x00, 0x0F, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorMalformedValueWithTracebackID("some-error", [16]byte{0xFF, 0xEE})},
			[]byte{0x00, 0x00, 0x09, 0xFF, 0xEE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x73, 0x6F, 0x6D, 0x65, 0x2D, 0x65, 0x72, 0x72, 0x6F, 0x72},
		},
		{
			ErrorAnswer{errs.NewErrorMismatchedChecksumWithTracebackID(0, 1, [16]byte{0x00, 0xAA})},
			[]byte{0x00, 0x00, 0x07, 0x00, 0xAA, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorMismatchedChecksumWithTracebackID(1, 1, [16]byte{0xBB, 0xCC, 0x01})},
			[]byte{0x00, 0x00, 0x07, 0xBB, 0xCC, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorNoHelloWithTracebackID([16]byte{0x01, 0x02, 0x03, 0x04})},
			[]byte{0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorNotFoundWithTracebackID([]byte("key"), [16]byte{0x03, 0x01, 0x02})},
			[]byte{0x00, 0x00, 0x0A, 0x03, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorNotFoundWithTracebackID([]byte(""), [16]byte{0xAB, 0xCD, 0xEF})},
			[]byte{0x00, 0x00, 0x0A, 0xAB, 0xCD, 0xEF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorProhibitedCompressionWithTracebackID(fields.S2, fields.Handshake, [16]byte{0xC0, 0xFF, 0xEE})},
			[]byte{0x00, 0x00, 0x0B, 0xC0, 0xFF, 0xEE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorProhibitedCompressionWithTracebackID(fields.None, fields.Handshake, [16]byte{0xC0, 0xFF, 0xEE})},
			[]byte{0x00, 0x00, 0x0B, 0xC0, 0xFF, 0xEE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorRequestsConflictWithTracebackID(0, [16]byte{0xFA, 0xCE, 0xB0, 0x0C})},
			[]byte{0x00, 0x00, 0x04, 0xFA, 0xCE, 0xB0, 0x0C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorRequestInterruptedWithTracebackID(0, [16]byte{0x0F, 0x0E, 0x0D})},
			[]byte{0x00, 0x00, 0x0C, 0x0F, 0x0E, 0x0D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorRestrictedRequestWithTracebackID([]byte("user"), []byte("1.1.1.1"), []byte("key"), fields.Read, 0, [16]byte{0x0C, 0x0B, 0x0A})},
			[]byte{0x00, 0x00, 0x11, 0x0C, 0x0B, 0x0A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnauthorizedWithTracebackID([]byte("1.1.1.1"), [16]byte{})},
			[]byte{0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnexpectedCommandInBatchWithTracebackID(fields.Answer, 0, [16]byte{})},
			[]byte{0x00, 0x00, 0x0E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnexpectedCommandInBatchWithTracebackID(fields.Read, 0, [16]byte{})},
			[]byte{0x00, 0x00, 0x0E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnexpectedCommandWithTracebackID(fields.Handshake, fields.Read, [16]byte{})},
			[]byte{0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
		},
		{
			ErrorAnswer{errs.NewErrorUnexpectedCommandWithTracebackID(fields.Handshake, fields.Handshake, [16]byte{})},
			[]byte{0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnsupportedCommandWithTracebackID(fields.Command(255), [16]byte{})},
			[]byte{0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnsupportedCommandWithTracebackID(fields.Handshake, [16]byte{})},
			[]byte{0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnsupportedCompressionWithTracebackID(fields.S2, [16]byte{})},
			[]byte{0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnsupportedCompressionWithTracebackID(fields.S2, [16]byte{})},
			[]byte{0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnsupportedVersionWithTracebackID(3, [16]byte{})},
			[]byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestErrorAnswer_Encode %v", tt.e),
			func(t *testing.T) {
				buf := buffer.Mock{}
				buf.Preallocate(tt.e.Size())

				tt.e.Encode(&buf)

				got := buf.Bytes()

				if len(got) != len(tt.want) {
					t.Fatalf("got %v want %v", len(got), len(tt.want))
				}

				if !bytes.Equal(got, tt.want) {
					t.Fatalf("got %q, want %q", got, tt.want)
				}
			},
		)
	}
}

func TestErrorAnswer_Command(t *testing.T) {
	tests := []struct {
		e    ErrorAnswer
		want fields.Command
	}{
		{ErrorAnswer{}, fields.Answer},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(1, 0)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(0, 1)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorBodyLimitIsExceeded(1, 0)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(0, 1)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorInternalError(errors.New(""))}, fields.Answer},
		{ErrorAnswer{errs.NewErrorInternalError(nil)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorInvalidRequestID(fields.Handshake, 1)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorInvalidRequestID(fields.Handshake, 0)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorMalformedValue("some-error")}, fields.Answer},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(0, 1)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(1, 1)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorNoHello()}, fields.Answer},
		{ErrorAnswer{errs.NewErrorNotFound([]byte("key"))}, fields.Answer},
		{ErrorAnswer{errs.NewErrorNotFound([]byte(""))}, fields.Answer},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(fields.S2, fields.Handshake)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(fields.None, fields.Handshake)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorRequestsConflict(0)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorRequestInterrupted(0)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorRestrictedRequest([]byte("user"), []byte("1.1.1.1"), []byte("key"), fields.Read, 0)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorUnauthorized([]byte("1.1.1.1"))}, fields.Answer},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(fields.Answer, 0)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(fields.Read, 0)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(fields.Handshake, fields.Read)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(fields.Handshake, fields.Handshake)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(fields.Command(255))}, fields.Answer},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(fields.Handshake)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(fields.S2)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(fields.None)}, fields.Answer},
		{ErrorAnswer{errs.NewErrorUnsupportedVersion(3)}, fields.Answer},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestErrorAnswer_Command %v", tt.e),
			func(t *testing.T) {
				if got := tt.e.Command(); got != tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestErrorAnswer_IsValid(t *testing.T) {
	tests := []struct {
		e    ErrorAnswer
		want bool
	}{
		{ErrorAnswer{}, true},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(1, 0)}, false},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(0, 1)}, true},
		{ErrorAnswer{errs.NewErrorBodyLimitIsExceeded(1, 0)}, false},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(0, 1)}, true},
		{ErrorAnswer{errs.NewErrorInternalError(errors.New(""))}, false},
		{ErrorAnswer{errs.NewErrorInternalError(nil)}, true},
		{ErrorAnswer{errs.NewErrorInvalidRequestID(fields.Handshake, 1)}, false},
		{ErrorAnswer{errs.NewErrorInvalidRequestID(fields.Handshake, 0)}, true},
		{ErrorAnswer{errs.NewErrorMalformedValue("some-error")}, false},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(0, 1)}, false},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(1, 1)}, true},
		{ErrorAnswer{errs.NewErrorNoHello()}, false},
		{ErrorAnswer{errs.NewErrorNotFound([]byte("key"))}, false},
		{ErrorAnswer{errs.NewErrorNotFound([]byte(""))}, true},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(fields.S2, fields.Handshake)}, false},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(fields.None, fields.Handshake)}, true},
		{ErrorAnswer{errs.NewErrorRequestsConflict(0)}, false},
		{ErrorAnswer{errs.NewErrorRequestInterrupted(0)}, false},
		{ErrorAnswer{errs.NewErrorRestrictedRequest([]byte("user"), []byte("1.1.1.1"), []byte("key"), fields.Read, 0)}, false},
		{ErrorAnswer{errs.NewErrorUnauthorized([]byte("1.1.1.1"))}, false},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(fields.Answer, 0)}, false},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(fields.Read, 0)}, true},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(fields.Handshake, fields.Read)}, false},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(fields.Handshake, fields.Handshake)}, true},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(fields.Command(255))}, false},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(fields.Handshake)}, true},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(fields.S2)}, false},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(fields.None)}, true},
		{ErrorAnswer{errs.NewErrorUnsupportedVersion(3)}, false},
	}
	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("TestErrorAnswer_IsValid %v", tt.e),
			func(t *testing.T) {
				if got := tt.e.IsValid(); got == nil == tt.want {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			},
		)
	}
}
