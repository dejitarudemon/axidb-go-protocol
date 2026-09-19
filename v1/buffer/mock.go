package buffer

import "encoding/binary"

var (
	_ Buffer = (*BufferMock)(nil)
)

type BufferMock struct {
	data []byte
}

func (b *BufferMock) Preallocate(size int) {
	b.data = make([]byte, 0, size)
}

func (b *BufferMock) Clean() {
	b.data = b.data[:0]
}

func (b *BufferMock) Bytes() []byte {
	return b.data
}

func (b *BufferMock) Append(v []byte) {
	b.data = append(b.data, v...)
}

func (b *BufferMock) AppendString(v string) {
	b.data = append(b.data, v...)
}

func (b *BufferMock) AppendUint8(v uint8) {
	b.data = append(b.data, byte(v))
}

func (b *BufferMock) AppendUint16(v uint16) {
	b.data = binary.BigEndian.AppendUint16(b.data, v)
}

func (b *BufferMock) AppendUint32(v uint32) {
	b.data = binary.BigEndian.AppendUint32(b.data, v)
}

func (b *BufferMock) AppendUint64(v uint64) {
	b.data = binary.BigEndian.AppendUint64(b.data, v)
}
