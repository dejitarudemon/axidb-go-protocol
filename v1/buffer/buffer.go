package buffer

type Appender interface {
	Append(v []byte)
	AppendUint8(v uint8)
	AppendUint16(v uint16)
	AppendUint32(v uint32)
	AppendUint64(v uint64)
	AppendString(v string)
}

type Buffer interface {
	Appender

	Preallocate(size int)
	Clean()
	Bytes() []byte
}
