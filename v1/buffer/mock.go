package buffer

import "encoding/binary"

var (
	_ Buffer = (*Mock)(nil)
)

type Mock struct {
	data []byte
}

func (m *Mock) Preallocate(size int) {
	m.data = make([]byte, 0, size)
}

func (m *Mock) Clean() {
	m.data = m.data[:0]
}

func (m *Mock) Bytes() []byte {
	return m.data
}

func (m *Mock) Append(v []byte) {
	m.data = append(m.data, v...)
}

func (m *Mock) AppendString(v string) {
	m.data = append(m.data, v...)
}

func (m *Mock) AppendUint8(v uint8) {
	m.data = append(m.data, byte(v))
}

func (m *Mock) AppendUint16(v uint16) {
	m.data = binary.BigEndian.AppendUint16(m.data, v)
}

func (m *Mock) AppendUint32(v uint32) {
	m.data = binary.BigEndian.AppendUint32(m.data, v)
}

func (m *Mock) AppendUint64(v uint64) {
	m.data = binary.BigEndian.AppendUint64(m.data, v)
}
