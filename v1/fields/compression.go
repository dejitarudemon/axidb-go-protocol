/*
package compression предназначен для представления кодов
сжатия (Compression) согласно спецификации протокола v1.

Использование:

	с = Code(1)
*/
package fields

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type Code предназначен для хранения
кода сжатия, его валидации и кодирования в сообщение.
*/
type Code uint8

/*
Константы, представляющие алгоритмы сжатия,
используемые в спецификации протокола v1.
*/
const (
	None Code = iota
	Zstd
	Lz4
)

/*
FieldSize представляет размер в байтах,
отведенный для хранения кода сжатия в сообщении.
*/
const FieldSize = 1

/*
func Encode предназначена для кодирования кода сжатия
в сообщении.

Принимааемые параметры:
  - buf buffer.Appender - буфер для хранения закодированного значения.
*/
func (c Code) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(c))
}

/*
func String предназначена для вывода человекочитаемого названия
алгоритма сжатия, представленного конкретным кодом.
*/
func (c Code) String() string {
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
func (c Code) IsValid() bool {
	return c <= Lz4
}

/*
func Size возвращает размер кода сжатия в байтах.
*/
func (c Code) Size() int {
	return FieldSize
}
