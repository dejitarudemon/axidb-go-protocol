/*
package compression предназначен для представления кодов
сжатия (Compression) согласно спецификации протокола v1

Использование:

	с = Compression(1)
*/
package compression

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/error/errs"
)

/*
type Compression предназначен для хранения
кода сжатия, его валидации и кодирования в сообщение
*/
type Compression uint8

/*
Константы, представляющие алгоритмы сжатия,
используемые в спецификации протокола v1
*/
const (
	None Compression = iota
	Zstd
	Lz4
)

/*
CompressionFieldSize представляет размер в байтах,
отведенный для хранения кода сжатия в сообщении
*/
const CompressionFieldSize = 1

/*
func Encode предназначена для кодирования кода сжатия
в сообщении.

Принимааемые параметры:
  - buf buffer.Appender - буфер для хранения закодированного значения
*/
func (c Compression) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(c))
}

/*
func String предназначена для вывода человекочитаемого названия
алгоритма сжатия, представленного конркетным кодом.
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

Возвращаемые ошибки:
 1. ErrorUnsupportedCompression - код находится вне пределов 0-2.
*/
func (c Compression) IsValid() error {
	if c > Lz4 {
		return errs.NewErrorUnsupportedCompression(uint8(c))
	}

	return nil
}

/*
func Size возвращает размер кода сжатия в байтах.
*/
func (c Compression) Size() int {
	return CompressionFieldSize
}
