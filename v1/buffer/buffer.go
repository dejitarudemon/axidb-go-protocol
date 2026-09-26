package buffer

// Appender writes encoded values into a buffer.
type Appender interface {
	// Append copies v into the buffer.
	Append(v []byte)

	// AppendUint8 appends v as a single byte.
	AppendUint8(v uint8)

	// AppendUint16 appends v in big-endian order.
	AppendUint16(v uint16)

	// AppendUint32 appends v in big-endian order.
	AppendUint32(v uint32)

	// AppendUint64 appends v in big-endian order.
	AppendUint64(v uint64)

	// AppendString appends the raw bytes of v.
	AppendString(v string)
}

// Buffer is an [Appender] that can reserve space, reset, and return its contents.
type Buffer interface {
	Appender

	// Preallocate ensures the buffer can hold at least size bytes without reallocating.
	Preallocate(size int)

	// Clean discards buffered bytes while keeping allocated capacity.
	Clean()

	// Raw returns the buffered bytes.
	// The slice aliases the buffer and is invalid after the next mutation.
	Raw() []byte

	// Bytes returns a copy of the buffered bytes.
	Bytes() []byte
}
