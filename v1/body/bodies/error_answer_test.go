package bodies

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compression"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
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
		{ErrorAnswer{errs.NewErrorInvalidRequestID(command.Handshake, 1)}, 19},
		{ErrorAnswer{errs.NewErrorInvalidRequestID(command.Handshake, 0)}, 19},
		{ErrorAnswer{errs.NewErrorMalformedValue("some-error")}, 29},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(0, 1)}, 19},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(1, 1)}, 19},
		{ErrorAnswer{errs.NewErrorNoHello()}, 19},
		{ErrorAnswer{errs.NewErrorNotFound([]byte("key"))}, 19},
		{ErrorAnswer{errs.NewErrorNotFound([]byte(""))}, 19},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(compression.Lz4, command.Handshake)}, 19},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(compression.None, command.Handshake)}, 19},
		{ErrorAnswer{errs.NewErrorRequestsConflict(0)}, 19},
		{ErrorAnswer{errs.NewErrorRequestInterrupted(0)}, 19},
		{ErrorAnswer{errs.NewErrorRestrictedRequest([]byte("user"), []byte("1.1.1.1"), []byte("key"), command.Read, 0)}, 19},
		{ErrorAnswer{errs.NewErrorUnauthorized([]byte("1.1.1.1"))}, 19},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(command.Answer, 0)}, 23},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(command.Read, 0)}, 23},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(command.Handshake, command.Read)}, 20},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(command.Handshake, command.Handshake)}, 20},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(command.Code(255))}, 19},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(command.Handshake)}, 19},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(compression.Lz4)}, 19},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(compression.None)}, 19},
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
			ErrorAnswer{errs.NewErrorInvalidRequestIDWithTracebackID(command.Handshake, 1, [16]byte{0x01})},
			[]byte{0x00, 0x00, 0x0F, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorInvalidRequestIDWithTracebackID(command.Handshake, 0, [16]byte{0x02, 0x03})},
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
			ErrorAnswer{errs.NewErrorProhibitedCompressionWithTracebackID(compression.Lz4, command.Handshake, [16]byte{0xC0, 0xFF, 0xEE})},
			[]byte{0x00, 0x00, 0x0B, 0xC0, 0xFF, 0xEE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorProhibitedCompressionWithTracebackID(compression.None, command.Handshake, [16]byte{0xC0, 0xFF, 0xEE})},
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
			ErrorAnswer{errs.NewErrorRestrictedRequestWithTracebackID([]byte("user"), []byte("1.1.1.1"), []byte("key"), command.Read, 0, [16]byte{0x0C, 0x0B, 0x0A})},
			[]byte{0x00, 0x00, 0x11, 0x0C, 0x0B, 0x0A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnauthorizedWithTracebackID([]byte("1.1.1.1"), [16]byte{})},
			[]byte{0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnexpectedCommandInBatchWithTracebackID(command.Answer, 0, [16]byte{})},
			[]byte{0x00, 0x00, 0x0E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnexpectedCommandInBatchWithTracebackID(command.Read, 0, [16]byte{})},
			[]byte{0x00, 0x00, 0x0E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnexpectedCommandWithTracebackID(command.Handshake, command.Read, [16]byte{})},
			[]byte{0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
		},
		{
			ErrorAnswer{errs.NewErrorUnexpectedCommandWithTracebackID(command.Handshake, command.Handshake, [16]byte{})},
			[]byte{0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnsupportedCommandWithTracebackID(command.Code(255), [16]byte{})},
			[]byte{0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnsupportedCommandWithTracebackID(command.Handshake, [16]byte{})},
			[]byte{0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnsupportedCompressionWithTracebackID(compression.Lz4, [16]byte{})},
			[]byte{0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			ErrorAnswer{errs.NewErrorUnsupportedCompressionWithTracebackID(compression.Lz4, [16]byte{})},
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
		want command.Code
	}{
		{ErrorAnswer{}, command.Answer},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(1, 0)}, command.Answer},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(0, 1)}, command.Answer},
		{ErrorAnswer{errs.NewErrorBodyLimitIsExceeded(1, 0)}, command.Answer},
		{ErrorAnswer{errs.NewErrorBatchLimitIsExceeded(0, 1)}, command.Answer},
		{ErrorAnswer{errs.NewErrorInternalError(errors.New(""))}, command.Answer},
		{ErrorAnswer{errs.NewErrorInternalError(nil)}, command.Answer},
		{ErrorAnswer{errs.NewErrorInvalidRequestID(command.Handshake, 1)}, command.Answer},
		{ErrorAnswer{errs.NewErrorInvalidRequestID(command.Handshake, 0)}, command.Answer},
		{ErrorAnswer{errs.NewErrorMalformedValue("some-error")}, command.Answer},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(0, 1)}, command.Answer},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(1, 1)}, command.Answer},
		{ErrorAnswer{errs.NewErrorNoHello()}, command.Answer},
		{ErrorAnswer{errs.NewErrorNotFound([]byte("key"))}, command.Answer},
		{ErrorAnswer{errs.NewErrorNotFound([]byte(""))}, command.Answer},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(compression.Lz4, command.Handshake)}, command.Answer},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(compression.None, command.Handshake)}, command.Answer},
		{ErrorAnswer{errs.NewErrorRequestsConflict(0)}, command.Answer},
		{ErrorAnswer{errs.NewErrorRequestInterrupted(0)}, command.Answer},
		{ErrorAnswer{errs.NewErrorRestrictedRequest([]byte("user"), []byte("1.1.1.1"), []byte("key"), command.Read, 0)}, command.Answer},
		{ErrorAnswer{errs.NewErrorUnauthorized([]byte("1.1.1.1"))}, command.Answer},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(command.Answer, 0)}, command.Answer},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(command.Read, 0)}, command.Answer},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(command.Handshake, command.Read)}, command.Answer},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(command.Handshake, command.Handshake)}, command.Answer},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(command.Code(255))}, command.Answer},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(command.Handshake)}, command.Answer},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(compression.Lz4)}, command.Answer},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(compression.None)}, command.Answer},
		{ErrorAnswer{errs.NewErrorUnsupportedVersion(3)}, command.Answer},
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
		{ErrorAnswer{errs.NewErrorInvalidRequestID(command.Handshake, 1)}, false},
		{ErrorAnswer{errs.NewErrorInvalidRequestID(command.Handshake, 0)}, true},
		{ErrorAnswer{errs.NewErrorMalformedValue("some-error")}, false},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(0, 1)}, false},
		{ErrorAnswer{errs.NewErrorMismatchedChecksum(1, 1)}, true},
		{ErrorAnswer{errs.NewErrorNoHello()}, false},
		{ErrorAnswer{errs.NewErrorNotFound([]byte("key"))}, false},
		{ErrorAnswer{errs.NewErrorNotFound([]byte(""))}, true},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(compression.Lz4, command.Handshake)}, false},
		{ErrorAnswer{errs.NewErrorProhibitedCompression(compression.None, command.Handshake)}, true},
		{ErrorAnswer{errs.NewErrorRequestsConflict(0)}, false},
		{ErrorAnswer{errs.NewErrorRequestInterrupted(0)}, false},
		{ErrorAnswer{errs.NewErrorRestrictedRequest([]byte("user"), []byte("1.1.1.1"), []byte("key"), command.Read, 0)}, false},
		{ErrorAnswer{errs.NewErrorUnauthorized([]byte("1.1.1.1"))}, false},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(command.Answer, 0)}, false},
		{ErrorAnswer{errs.NewErrorUnexpectedCommandInBatch(command.Read, 0)}, true},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(command.Handshake, command.Read)}, false},
		{ErrorAnswer{errs.NewErrorUnexpectedCommand(command.Handshake, command.Handshake)}, true},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(command.Code(255))}, false},
		{ErrorAnswer{errs.NewErrorUnsupportedCommand(command.Handshake)}, true},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(compression.Lz4)}, false},
		{ErrorAnswer{errs.NewErrorUnsupportedCompression(compression.None)}, true},
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
