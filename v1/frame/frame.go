package frame

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

type DefaultBufferForBodyCompression = buffer.Slice

var MagicBytes = []byte{0x0A, 0xDB}

const (
	HeadersLen = 10
)

type Frame struct {
	RequestID fields.RequestID
	Body      body.Body
}

func (f Frame) Size() int {
	size := f.RequestID.Size() + body.BodyLenFieldSize + len(MagicBytes) + f.Version().Size() + fields.CompressionFieldSize + fields.ChecksumFieldSize
	if f.Body == nil {
		return fields.CommandFieldSize + size
	}

	return f.Body.Command().Size() + f.Body.Size() + size
}

func (f Frame) IsValid() error {
	if f.Body == nil {
		return err.NewValidationError(
			"nil Body",
			"target", "Frame",
		)
	}

	return f.Body.IsValid()
}

func (f Frame) Version() fields.Version {
	return fields.CurrentVersion
}

func (f Frame) encodeWithCompression(buf buffer.Buffer, compressor compressor.Compressor) error {
	temp := DefaultBufferForBodyCompression{}
	temp.Preallocate(f.Body.Size())

	f.Body.Encode(&temp)

	compressed, err := compressor.Compress(temp.Bytes())
	if err != nil {
		return err
	}

	f.Body.Command().Encode(buf)
	f.RequestID.Encode(buf)
	compressor.Code().Encode(buf)
	buf.AppendUint32(uint32(len(compressed)))
	buf.Append(compressed)

	fields.NewChecksum(buf.Bytes()).Encode(buf)

	return nil
}

func (f Frame) encodeWithoutCompression(buf buffer.Buffer) {
	f.Body.Command().Encode(buf)

	f.RequestID.Encode(buf)
	fields.None.Encode(buf)

	buf.AppendUint32(uint32(f.Body.Size()))
	f.Body.Encode(buf)

	fields.NewChecksum(buf.Bytes()).Encode(buf)
}

func (f Frame) Encode(buf buffer.Buffer, compressor compressor.Compressor) error {
	if f.Body == nil {
		return err.NewValidationError(
			"nil Body",
			"target", "Frame",
		)
	}

	buf.Append(MagicBytes)
	f.Version().Encode(buf)

	if compressor == nil {
		f.encodeWithoutCompression(buf)
		return nil
	}

	return f.encodeWithCompression(buf, compressor)
}
