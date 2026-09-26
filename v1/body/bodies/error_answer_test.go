package bodies

import (
	"errors"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func TestErrorAnswer(t *testing.T) {
	id := fields.TracebackID{15: 0x01}

	tests := []struct {
		name string
		err  err.ProtocolError
	}{
		{"nil error", nil},
		{"no hello", errs.NewErrorNoHelloWithTracebackID(id)},
		{"unsupported version", errs.NewErrorUnsupportedVersionWithTracebackID(3, id)},
		{"unexpected command", errs.NewErrorUnexpectedCommandWithTracebackID(fields.Handshake, fields.Read, id)},
		{"unsupported command", errs.NewErrorUnsupportedCommandWithTracebackID(fields.Command(255), id)},
		{"requests conflict", errs.NewErrorRequestsConflictWithTracebackID(0, id)},
		{"unsupported compression", errs.NewErrorUnsupportedCompressionWithTracebackID(fields.S2, id)},
		{"body limit is exceeded", errs.NewErrorBodyLimitIsExceededWithTracebackID(1, 0, id)},
		{"mismatched checksum", errs.NewErrorMismatchedChecksumWithTracebackID(0, 1, id)},
		{"internal error", errs.NewErrorInternalErrorWithTracebackID(errors.New(""), id)},
		{"malformed value", errs.NewErrorMalformedValueWithTracebackID("some-error", id)},
		{"not found", errs.NewErrorNotFoundWithTracebackID(fields.Key("key"), id)},
		{"prohibited compression", errs.NewErrorProhibitedCompressionWithTracebackID(fields.S2, fields.Handshake, id)},
		{"request interrupted", errs.NewErrorRequestInterruptedWithTracebackID(0, id)},
		{"batch limit is exceeded", errs.NewErrorBatchLimitIsExceededWithTracebackID(1, 0, id)},
		{"unexpected command in batch", errs.NewErrorUnexpectedCommandInBatchWithTracebackID(fields.Answer, 0, id)},
		{"invalid request id", errs.NewErrorInvalidRequestIDWithTracebackID(fields.Handshake, 1, id)},
		{"unauthorized", errs.NewErrorUnauthorizedWithTracebackID([]byte("1.1.1.1"), id)},
		{"restricted request", errs.NewErrorRestrictedRequestWithTracebackID(fields.Key("user"), fields.Key("1.1.1.1"), fields.Key("key"), fields.Read, 0, id)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := ErrorAnswer{Err: tt.err}
			want := []byte{0x00}
			if tt.err != nil {
				want = append(want, testutil.Encode(t, tt.err)...)
			}

			assertAnswer(t, b, fields.Answer, want, tt.err == nil)
		})
	}
}
