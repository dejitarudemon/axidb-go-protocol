package fields

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type Compression предназначен для хранения
кода сжатия, его валидации и кодирования в сообщение.
*/
type Compression uint8

/*
Константы, представляющие алгоритмы сжатия,
используемые в спецификации протокола v1.
*/
const (
	None Compression = iota
	Zstd
	Lz4
)

/*
CompressionFieldSize представляет размер в байтах,
отведенный для хранения кода сжатия в сообщении.
*/
const CompressionFieldSize = 1

/*
func Encode предназначена для кодирования кода сжатия
в сообщении.

Принимааемые параметры:
  - buf buffer.Appender - буфер для хранения закодированного значения.
*/
func (c Compression) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(c))
}

/*
func String предназначена для вывода человекочитаемого названия
алгоритма сжатия, представленного конкретным кодом.
*/
func (c Compression) String() string {
	switch c {
	case None:
		return "None"
	case Zstd:
		return "Zstd"
	case Lz4:
		return "Lz4"
	}

	return fmt.Sprintf("Unknown (%d)", c)
}

/*
func IsValid предназначена для проверки кода сжатия.
Проверки:
 1. Код находится в пределах 0-2.
*/
func (c Compression) IsValid() bool {
	return c <= Lz4
}

/*
func Size возвращает размер кода сжатия в байтах.
*/
func (c Compression) Size() int {
	return CompressionFieldSize
}
