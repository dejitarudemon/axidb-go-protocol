package buffer

import "encoding/binary"

var (
	_ Buffer = (*Slice)(nil)
)

type Slice struct {
	data []byte
}

func (s *Slice) Preallocate(size int) {
	if cap(s.data) < size {
		s.data = make([]byte, 0, size)
	}
}

func (s *Slice) Clean() {
	s.data = s.data[:0]
}

func (s *Slice) Bytes() []byte {
	return s.data[:len(s.data)]
}

func (s *Slice) Append(v []byte) {
	s.data = append(s.data, v...)
}

func (s *Slice) AppendString(v string) {
	s.data = append(s.data, v...)
}

func (s *Slice) AppendUint8(v uint8) {
	s.data = append(s.data, byte(v))
}

func (s *Slice) AppendUint16(v uint16) {
	s.data = binary.BigEndian.AppendUint16(s.data, v)
}

func (s *Slice) AppendUint32(v uint32) {
	s.data = binary.BigEndian.AppendUint32(s.data, v)
}

func (s *Slice) AppendUint64(v uint64) {
	s.data = binary.BigEndian.AppendUint64(s.data, v)
}
