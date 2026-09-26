package decoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/internal/testutil"
)

func cat(parts ...[]byte) []byte {
	return bytes.Join(parts, nil)
}

func u16(v uint16) []byte {
	return binary.BigEndian.AppendUint16(nil, v)
}

func u32(v uint32) []byte {
	return binary.BigEndian.AppendUint32(nil, v)
}

func u64(v uint64) []byte {
	return binary.BigEndian.AppendUint64(nil, v)
}

func encodeBody(t testing.TB, e testutil.Encoder) []byte {
	t.Helper()
	return testutil.Encode(t, e)
}

func assertMalformed(t *testing.T, e error) {
	t.Helper()

	var mv errs.ErrorMalformedValue
	if !errors.As(e, &mv) {
		t.Fatalf("got %v, want ErrorMalformedValue", e)
	}
}

func assertDecoded(t *testing.T, got testutil.Encoder, e error, c *cursor, want testutil.Encoder) {
	t.Helper()

	if e != nil {
		t.Fatalf("decode = %v", e)
	}

	if e := c.expectEnd(); e != nil {
		t.Fatalf("trailing data: %v", e)
	}

	testutil.AssertSameEncoding(t, got, want)
}
