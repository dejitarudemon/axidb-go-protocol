package fields

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type Type предназначен для хранения
кода типа, его валидации и кодирования в сообщение.
*/
type Type uint8

/*
FieldSize представляет размер в байтах,
отведенный для хранения кода типа в сообщении.
*/
const TypeFieldSize = 1

/*
Константы, представляющие типов данных,
используемые в спецификации протокола v1.
*/
const (
	Bytes Type = iota
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
func (t Type) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(t))
}

/*
func String предназначена для вывода человекочитаемого названия
типа данных, представленного конкретным кодом.
*/
func (t Type) String() string {
	switch t {
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

	return fmt.Sprintf("Unknown (%d)", t)
}

/*
func IsValid предназначена для проверки кода типа данных.
Проверки:
 1. Код находится в пределах 0-7.
*/
func (t Type) IsValid() bool {
	return t <= JSON
}

/*
func Size возвращает размер кода типа данных в байтах.
*/
func (t Type) Size() int {
	return TypeFieldSize
}
