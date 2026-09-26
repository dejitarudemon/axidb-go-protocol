package decoder

import (
	"encoding/binary"
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
)

// cursor reads big-endian fields from a byte slice and reports a malformed value
// instead of panicking when the slice is too short.
type cursor struct {
	buf []byte
	pos int
}

// newCursor returns a cursor positioned at the start of buf.
func newCursor(buf []byte) *cursor {
	return &cursor{buf: buf}
}

// remaining returns the number of unread bytes.
func (c *cursor) remaining() int {
	return len(c.buf) - c.pos
}

// bytes returns the next n bytes and advances past them.
// The returned slice aliases the underlying buffer.
func (c *cursor) bytes(n int) ([]byte, error) {
	if n < 0 || n > c.remaining() {
		return nil, errs.NewErrorMalformedValue(
			fmt.Sprintf("unexpected end of body: need %v bytes at offset %v, got %v", n, c.pos, c.remaining()),
		)
	}

	b := c.buf[c.pos : c.pos+n]
	c.pos += n

	return b, nil
}

// rest returns all unread bytes and advances to the end.
func (c *cursor) rest() []byte {
	b := c.buf[c.pos:]
	c.pos = len(c.buf)

	return b
}

// sub returns a cursor over the next n bytes and advances past them.
func (c *cursor) sub(n int) (*cursor, error) {
	b, e := c.bytes(n)
	if e != nil {
		return nil, e
	}

	return newCursor(b), nil
}

// uint8 reads one byte.
func (c *cursor) uint8() (uint8, error) {
	b, e := c.bytes(1)
	if e != nil {
		return 0, e
	}

	return b[0], nil
}

// uint16 reads a big-endian uint16.
func (c *cursor) uint16() (uint16, error) {
	b, e := c.bytes(2)
	if e != nil {
		return 0, e
	}

	return binary.BigEndian.Uint16(b), nil
}

// uint32 reads a big-endian uint32.
func (c *cursor) uint32() (uint32, error) {
	b, e := c.bytes(4)
	if e != nil {
		return 0, e
	}

	return binary.BigEndian.Uint32(b), nil
}

// uint64 reads a big-endian uint64.
func (c *cursor) uint64() (uint64, error) {
	b, e := c.bytes(8)
	if e != nil {
		return 0, e
	}

	return binary.BigEndian.Uint64(b), nil
}

// length reads a big-endian uint32 length prefix as an int.
func (c *cursor) length() (int, error) {
	n, e := c.uint32()
	return int(n), e
}

// expectEnd reports trailing bytes when the cursor has not consumed its whole buffer.
func (c *cursor) expectEnd() error {
	if c.remaining() != 0 {
		return errs.NewErrorMalformedValue(
			fmt.Sprintf("unexpected trailing data: %v of %v bytes left unread", c.remaining(), len(c.buf)),
		)
	}

	return nil
}
