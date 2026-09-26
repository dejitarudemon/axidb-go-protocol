package decoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
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

func encodeBody(b body.Body) []byte {
	buf := buffer.Slice{}
	buf.Preallocate(b.Size())
	b.Encode(&buf)
	return buf.Bytes()
}

func assertMalformed(t *testing.T, e error) {
	t.Helper()

	var mv errs.ErrorMalformedValue
	if !errors.As(e, &mv) {
		t.Fatalf("expected ErrorMalformedValue, got %v", e)
	}
}

func TestCursor_Reads(t *testing.T) {
	c := newCursor(cat(
		[]byte{0x01},
		u16(0x0203),
		u32(0x04050607),
		u64(0x08090A0B0C0D0E0F),
		u32(2),
		[]byte{0xAA, 0xBB, 0xCC},
	))

	if v, e := c.uint8(); e != nil || v != 0x01 {
		t.Fatalf("uint8: got %v, %v", v, e)
	}
	if v, e := c.uint16(); e != nil || v != 0x0203 {
		t.Fatalf("uint16: got %v, %v", v, e)
	}
	if v, e := c.uint32(); e != nil || v != 0x04050607 {
		t.Fatalf("uint32: got %v, %v", v, e)
	}
	if v, e := c.uint64(); e != nil || v != 0x08090A0B0C0D0E0F {
		t.Fatalf("uint64: got %v, %v", v, e)
	}
	if v, e := c.length(); e != nil || v != 2 {
		t.Fatalf("length: got %v, %v", v, e)
	}
	if v, e := c.bytes(1); e != nil || !bytes.Equal(v, []byte{0xAA}) {
		t.Fatalf("bytes: got %v, %v", v, e)
	}
	if e := c.expectEnd(); e == nil {
		t.Fatal("expectEnd: expected err with unread bytes")
	}
	if v := c.rest(); !bytes.Equal(v, []byte{0xBB, 0xCC}) {
		t.Fatalf("rest: got %v", v)
	}
	if c.remaining() != 0 {
		t.Fatalf("remaining: got %v", c.remaining())
	}
	if e := c.expectEnd(); e != nil {
		t.Fatalf("expectEnd: got %v", e)
	}
	if v := c.rest(); len(v) != 0 {
		t.Fatalf("rest at end: got %v", v)
	}
}

func TestCursor_ShortReads(t *testing.T) {
	tests := []struct {
		name string
		read func(c *cursor) error
	}{
		{"uint8", func(c *cursor) error { _, e := c.uint8(); return e }},
		{"uint16", func(c *cursor) error { _, e := c.uint16(); return e }},
		{"uint32", func(c *cursor) error { _, e := c.uint32(); return e }},
		{"uint64", func(c *cursor) error { _, e := c.uint64(); return e }},
		{"length", func(c *cursor) error { _, e := c.length(); return e }},
		{"bytes", func(c *cursor) error { _, e := c.bytes(2); return e }},
		{"bytes negative", func(c *cursor) error { _, e := c.bytes(-1); return e }},
		{"sub", func(c *cursor) error { _, e := c.sub(2); return e }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCursor(nil)
			assertMalformed(t, tt.read(c))

			if c.pos != 0 {
				t.Errorf("failed read moved cursor to %v", c.pos)
			}
		})
	}
}

func TestCursor_Sub(t *testing.T) {
	c := newCursor([]byte{0x01, 0x02, 0x03})

	sub, e := c.sub(2)
	if e != nil {
		t.Fatalf("sub: got %v", e)
	}

	if !bytes.Equal(sub.rest(), []byte{0x01, 0x02}) {
		t.Error("sub: unexpected content")
	}

	if c.remaining() != 1 {
		t.Errorf("sub: parent remaining %v, want 1", c.remaining())
	}
}
