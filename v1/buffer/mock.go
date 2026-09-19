package buffer

import "encoding/binary"

var (
	_ Buffer = (*BufferMock)(nil)
)

type BufferMock []byte

func (b BufferMock) Preallocate(size int) {
	b = make(BufferMock, 0, size)
}

func (b BufferMock) Clean() {
	b = b[:0]
}

func (b BufferMock) Bytes() []byte {
	return b
}

func (a BufferMock) Append(v []byte) {
	a = append(a, v...)
}

func (a BufferMock) AppendString(v string) {
	a = append(a, v...)
}

func (a BufferMock) AppendUint8(v uint8) {
	a = append(a, byte(v))
}

func (a BufferMock) AppendUint16(v uint16) {
	a = binary.BigEndian.AppendUint16(a, v)
}

func (a BufferMock) AppendUint32(v uint32) {
	a = binary.BigEndian.AppendUint32(a, v)
}

func (a BufferMock) AppendUint64(v uint64) {
	a = binary.BigEndian.AppendUint64(a, v)
}
