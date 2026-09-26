package decoder

import (
	"bytes"
	"testing"
)

func TestCursor_Reads(t *testing.T) {
	c := newCursor(cat(
		[]byte{0x01},
		u16(0x0203),
		u32(0x04050607),
		u64(0x08090A0B0C0D0E0F),
		u32(2),
		[]byte{0xAA, 0xBB, 0xCC},
	))

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"uint8", mustRead(t, c.uint8), uint8(0x01)},
		{"uint16", mustRead(t, c.uint16), uint16(0x0203)},
		{"uint32", mustRead(t, c.uint32), uint32(0x04050607)},
		{"uint64", mustRead(t, c.uint64), uint64(0x08090A0B0C0D0E0F)},
		{"length", mustRead(t, c.length), 2},
		{"bytes", mustRead(t, func() ([]byte, error) { return c.bytes(1) }), []byte{0xAA}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch want := tt.want.(type) {
			case []byte:
				got, ok := tt.got.([]byte)
				if !ok || !bytes.Equal(got, want) {
					t.Errorf("got %v, want %v", tt.got, want)
				}
			default:
				if tt.got != tt.want {
					t.Errorf("got %v, want %v", tt.got, tt.want)
				}
			}
		})
	}

	if e := c.expectEnd(); e == nil {
		t.Fatal("expectEnd() = nil, want error")
	}

	if got, want := c.rest(), []byte{0xBB, 0xCC}; !bytes.Equal(got, want) {
		t.Errorf("rest() = % X, want % X", got, want)
	}

	if c.remaining() != 0 {
		t.Errorf("remaining() = %v, want 0", c.remaining())
	}

	if e := c.expectEnd(); e != nil {
		t.Errorf("expectEnd() = %v", e)
	}

	if got := c.rest(); len(got) != 0 {
		t.Errorf("rest() at end = % X, want empty", got)
	}
}

func mustRead[T any](t *testing.T, read func() (T, error)) T {
	t.Helper()

	got, e := read()
	if e != nil {
		t.Fatalf("read = %v", e)
	}

	return got
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
		t.Fatalf("sub() = %v", e)
	}

	if got, want := sub.rest(), []byte{0x01, 0x02}; !bytes.Equal(got, want) {
		t.Errorf("sub rest() = % X, want % X", got, want)
	}

	if c.remaining() != 1 {
		t.Errorf("parent remaining() = %v, want 1", c.remaining())
	}
}
