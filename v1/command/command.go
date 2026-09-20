/*
package command предназначен для представления кодов
команд (Command) согласно спецификации протокола v1.

Использование:

	с = Code(1)
*/

package command

import (
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
)

/*
type Code предназначен для хранения
кода команды, его валидации и кодирования в сообщение.
*/
type Code uint8

/*
FieldSize представляет размер в байтах,
отведенный для хранения кода команды в сообщении.
*/
const FieldSize = 1

/*
Константы, представляющие команды,
используемые в спецификации протокола v1.
*/
const (
	Handshake Code = iota
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
func (c Code) Encode(buf buffer.Appender) {
	buf.AppendUint8(uint8(c))
}

/*
func String предназначена для вывода человекочитаемого названия
команды, представленного конкретным кодом.
*/
func (c Code) String() string {
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

Возвращаемые ошибки:
 1. ErrorUnsupportedCommand - код находится вне пределов 0-6.
*/
func (c Code) IsValid() error {
	if c > Ping {
		return errs.NewErrorUnsupportedCommand(uint8(c))
	}

	return nil
}

/*
func Size возвращает размер кода команды в байтах.
*/
func (c Code) Size() int {
	return FieldSize
}
