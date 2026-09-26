package values

import (
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

// assertValue checks the methods shared by every value.
func assertValue(t *testing.T, v value.V, typ fields.Type, want []byte, wantErr bool) {
	t.Helper()

	if got := v.Type(); got != typ {
		t.Errorf("Type() = %v, want %v", got, typ)
	}

	testutil.AssertEncoded(t, v, want)
	testutil.AssertErr(t, v.IsValid(), wantErr)
}
