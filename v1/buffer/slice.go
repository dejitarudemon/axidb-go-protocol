package buffer

import (
	"encoding/binary"
	"sync"
)

var (
	_ Buffer = (*Slice)(nil)
)

type Slice struct {
	data []byte
	mx   sync.Mutex
}

func (s *Slice) Preallocate(size int) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if cap(s.data) < size {
		s.data = make([]byte, 0, size)
	}
}

func (s *Slice) Clean() {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.data = s.data[:0]
}

func (s *Slice) Bytes() []byte {
	s.mx.Lock()
	defer s.mx.Unlock()

	data := make([]byte, len(s.data))
	copy(data, s.data)

	return data
}

func (s *Slice) Append(v []byte) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.data = append(s.data, v...)
}

func (s *Slice) AppendString(v string) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.data = append(s.data, v...)
}

func (s *Slice) AppendUint8(v uint8) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.data = append(s.data, byte(v))
}

func (s *Slice) AppendUint16(v uint16) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.data = binary.BigEndian.AppendUint16(s.data, v)
}

func (s *Slice) AppendUint32(v uint32) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.data = binary.BigEndian.AppendUint32(s.data, v)
}

func (s *Slice) AppendUint64(v uint64) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.data = binary.BigEndian.AppendUint64(s.data, v)
}
