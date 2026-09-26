package bodies

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

// assertBody checks the methods shared by every body.
func assertBody(t *testing.T, b body.Body, command fields.Command, want []byte, wantErr bool) {
	t.Helper()

	if got := b.Command(); got != command {
		t.Errorf("Command() = %v, want %v", got, command)
	}

	testutil.AssertEncoded(t, b, want)
	testutil.AssertErr(t, b.IsValid(), wantErr)
}

// assertAnswer checks the methods shared by every answer body.
func assertAnswer(t *testing.T, a body.Answer, responseTo fields.Command, want []byte, wantErr bool) {
	t.Helper()

	if got := a.IsResponseTo(); got != responseTo {
		t.Errorf("IsResponseTo() = %v, want %v", got, responseTo)
	}

	assertBody(t, a, fields.Answer, want, wantErr)
}
