package types

/*
package package types предназначен для представления кодов
типов данных (Value Type) согласно спецификации протокола v1.

Использование:

	с = Code(1)
*/

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
)

/*
type Code предназначен для хранения
кода типа, его валидации и кодирования в сообщение.
*/
type Code uint8

/*
FieldSize представляет размер в байтах,
отведенный для хранения кода типа в сообщении.
*/
const FieldSize = 1

/*
Константы, представляющие типов данных,
используемые в спецификации протокола v1.
*/
const (
	Bytes Code = iota
	TypedArray
	UntypedArray
	Int
	Uint
	Float
	String
	JSON
)

/*
func Encode предназначена для кодирования кода типа данных
в сообщении.

Принимааемые параметры:
  - buf buffer.Appender - буфер для хранения закодированного значения.
*/
func (c Code) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(c))
}

/*
func String предназначена для вывода человекочитаемого названия
типа данных, представленного конкретным кодом.
*/
func (c Code) String() string {
	switch c {
	case Bytes:
		return "Bytes"
	case TypedArray:
		return "Typed Array"
	case UntypedArray:
		return "Untyped Array"
	case Int:
		return "Int"
	case Uint:
		return "Uint"
	case Float:
		return "Float"
	case String:
		return "String"
	case JSON:
		return "JSON"
	}

	return fmt.Sprintf("Unknown (%d)", c)
}

/*
func IsValid предназначена для проверки кода типа данных.
Проверки:
 1. Код находится в пределах 0-7.

Возвращаемые ошибки:
 1. ErrorUnsupportedCommand - код находится вне пределов 0-7.
*/
func (c Code) IsValid() error {
	if c > JSON {
		return errs.NewErrorUnsupportedCommand(uint8(c))
	}

	return nil
}

/*
func Size возвращает размер кода типа данных в байтах.
*/
func (c Code) Size() int {
	return FieldSize
}
