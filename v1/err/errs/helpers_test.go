package errs

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

var (
	tracebackIDOne = fields.TracebackID{15: 0x01}
	tracebackIDMax = fields.TracebackID{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}
)

// encoded builds the expected wire form of an error: code, traceback ID, then details.
func encoded(code uint16, id fields.TracebackID, details ...byte) []byte {
	return bytes.Join([][]byte{{byte(code >> 8), byte(code)}, id[:], details}, nil)
}

// assertProtocolError checks the methods shared by every protocol error.
func assertProtocolError(t *testing.T, e err.ProtocolError, code fields.Error, id fields.TracebackID, want []byte, wantErr bool) {
	t.Helper()

	if got := e.Code(); got != code {
		t.Errorf("Code() = %v, want %v", got, code)
	}

	if got := e.TracebackID(); got != id {
		t.Errorf("TracebackID() = %v, want %v", got, id)
	}

	if msg := e.Error(); !strings.Contains(msg, code.String()) {
		t.Errorf("Error() = %q, want it to mention %q", msg, code)
	}

	testutil.AssertEncoded(t, e, want)

	v, ok := e.(interface{ IsValid() error })
	if !ok {
		t.Fatalf("%T has no IsValid method", e)
	}

	testutil.AssertErr(t, v.IsValid(), wantErr)
}

func TestNewErrors_GenerateTracebackID(t *testing.T) {
	tests := []struct {
		name string
		e    err.ProtocolError
	}{
		{"NoHello", NewErrorNoHello()},
		{"UnsupportedVersion", NewErrorUnsupportedVersion(2)},
		{"UnexpectedCommand", NewErrorUnexpectedCommand(fields.Read, fields.Handshake)},
		{"UnsupportedCommand", NewErrorUnsupportedCommand(fields.Command(10))},
		{"RequestsConflict", NewErrorRequestsConflict(1)},
		{"UnsupportedCompression", NewErrorUnsupportedCompression(fields.Compression(3))},
		{"BodyLimitIsExceeded", NewErrorBodyLimitIsExceeded(2, 1)},
		{"MismatchedChecksum", NewErrorMismatchedChecksum(1, 2)},
		{"InternalError", NewErrorInternalError(nil)},
		{"MalformedValue", NewErrorMalformedValue("msg")},
		{"NotFound", NewErrorNotFound(fields.Key("key"))},
		{"ProhibitedCompression", NewErrorProhibitedCompression(fields.Zstd, fields.Ping)},
		{"RequestInterrupted", NewErrorRequestInterrupted(1)},
		{"BatchLimitIsExceeded", NewErrorBatchLimitIsExceeded(2, 1)},
		{"UnexpectedCommandInBatch", NewErrorUnexpectedCommandInBatch(fields.Ping, 1)},
		{"InvalidRequestID", NewErrorInvalidRequestID(fields.Read, 0)},
		{"Unauthorized", NewErrorUnauthorized([]byte("source"))},
		{"RestrictedRequest", NewErrorRestrictedRequest(fields.Key("id"), fields.Key("source"), fields.Key("key"), fields.Read, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.e.TracebackID() == (fields.TracebackID{}) {
				t.Error("TracebackID() is zero, want a generated one")
			}
		})
	}
}
