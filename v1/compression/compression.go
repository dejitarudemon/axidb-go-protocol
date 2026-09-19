package compression

type Compression uint8

const (
	None Compression = iota
	Zstd
	Lz4
)

const CompressionFieldSize = 1
