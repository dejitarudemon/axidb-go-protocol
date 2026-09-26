// Package testutil holds assertions shared by protocol v1 tests.
//
// It depends only on [buffer], so any v1 package except buffer itself can use it
// from internal tests without an import cycle.
package testutil

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

// Encoder is a protocol value with a wire encoding.
type Encoder interface {
	Size() int
	Encode(buf buffer.Appender)
}

// Encode returns the wire encoding of e and fails t when its length disagrees with e.Size().
func Encode(t testing.TB, e Encoder) []byte {
	t.Helper()

	buf := buffer.Slice{}
	buf.Preallocate(e.Size())
	e.Encode(&buf)
	got := buf.Bytes()

	if len(got) != e.Size() {
		t.Errorf("Size() = %v, but Encode() wrote %v bytes", e.Size(), len(got))
	}

	return got
}

// AssertEncoded fails t when e does not encode to want.
func AssertEncoded(t testing.TB, e Encoder, want []byte) {
	t.Helper()

	if got := Encode(t, e); !bytes.Equal(got, want) {
		t.Errorf("Encode() = % X, want % X", got, want)
	}
}

// AssertSameEncoding fails t when got and want differ in type or wire encoding.
// It compares what is actually sent, so fields that are not encoded are ignored.
func AssertSameEncoding(t testing.TB, got, want Encoder) {
	t.Helper()

	if got == nil || want == nil {
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		return
	}

	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("got %T, want %T", got, want)
		return
	}

	if g, w := Encode(t, got), Encode(t, want); !bytes.Equal(g, w) {
		t.Errorf("encoding mismatch for %T:\ngot  % X\nwant % X", got, g, w)
	}
}

// AssertErr fails t when the presence of err does not match wantErr.
func AssertErr(t testing.TB, err error, wantErr bool) {
	t.Helper()

	if (err != nil) != wantErr {
		t.Errorf("got err %v, want err: %v", err, wantErr)
	}
}