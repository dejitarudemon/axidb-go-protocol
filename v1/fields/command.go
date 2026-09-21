/*
package field предназначен для представления различных кодов и полей
согласно спецификации протокола v1.
*/

package fields

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

/*
type Command предназначен для хранения
кода команды, его валидации и кодирования в сообщение.
*/
type Command uint8

/*
FieldSize представляет размер в байтах,
отведенный для хранения кода команды в сообщении.
*/
const CommandFieldSize = 1

/*
Константы, представляющие команды,
используемые в спецификации протокола v1.
*/
const (
	Handshake Command = iota
	Answer
	Read
	Write
	Delete
	Batch
	Ping
)

/*
func Encode предназначена для кодирования кода команды
в сообщении.

Принимааемые параметры:
  - buf buffer.Appender - буфер для хранения закодированного значения.
*/
func (c Command) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(c))
}

/*
func String предназначена для вывода человекочитаемого названия
команды, представленного конкретным кодом.
*/
func (c Command) String() string {
	switch c {
	case Handshake:
		return "Handshake"
	case Answer:
		return "Answer"
	case Read:
		return "Read"
	case Write:
		return "Write"
	case Delete:
		return "Delete"
	case Batch:
		return "Batch"
	case Ping:
		return "Ping"
	}

	return fmt.Sprintf("Unknown (%d)", c)
}

/*
func IsValid предназначена для проверки кода команды.
Проверки:
 1. Код находится в пределах 0-6.
*/
func (c Command) IsValid() bool {
	return c <= Ping
}

/*
func Size возвращает размер кода команды в байтах.
*/
func (c Command) Size() int {
	return CommandFieldSize
}
