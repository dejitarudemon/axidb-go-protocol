package buffer

import (
	"encoding/binary"
)

var (
	_ Buffer = (*Slice)(nil)
)

// Slice is a [Buffer] backed by a byte slice.
// It is not safe for concurrent use. Integer values are appended in big-endian order.
type Slice struct {
	data []byte
}

// Preallocate ensures the slice can hold at least size bytes.
// If the current capacity is smaller, the buffer is replaced with an empty slice of that capacity.
func (s *Slice) Preallocate(size int) {
	if cap(s.data) < size {
		s.data = make([]byte, 0, size)
	}
}

// Clean discards buffered bytes while keeping allocated capacity.
func (s *Slice) Clean() {
	s.data = s.data[:0]
}

// Raw returns the buffered bytes.
// The slice aliases the buffer and is invalid after the next mutation.
func (s *Slice) Raw() []byte {
	return s.data
}

// Bytes returns a copy of the buffered bytes.
func (s *Slice) Bytes() []byte {
	data := make([]byte, len(s.data))
	copy(data, s.data)
	return data
}

// Append copies v into the buffer.
func (s *Slice) Append(v []byte) {
	s.data = append(s.data, v...)
}

// AppendString appends the raw bytes of v.
func (s *Slice) AppendString(v string) {
	s.data = append(s.data, v...)
}

// AppendUint8 appends v as a single byte.
func (s *Slice) AppendUint8(v uint8) {
	s.data = append(s.data, byte(v))
}

// AppendUint16 appends v in big-endian order.
func (s *Slice) AppendUint16(v uint16) {
	s.data = binary.BigEndian.AppendUint16(s.data, v)
}

// AppendUint32 appends v in big-endian order.
func (s *Slice) AppendUint32(v uint32) {
	s.data = binary.BigEndian.AppendUint32(s.data, v)
}

// AppendUint64 appends v in big-endian order.
func (s *Slice) AppendUint64(v uint64) {
	s.data = binary.BigEndian.AppendUint64(s.data, v)
}
