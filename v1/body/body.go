/*
package body предназначен для представления различных
тел сообщений (Body), представленных в спецификации протокола v1.
*/
package body

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/command"
)

type Body interface {
	/*
		Size возвращает размер тела в байтах.
	*/
	Size() int

	/*
		Encode предназначена для кодирования кода тела сообщения
		в сообщении.

		Принимааемые параметры:
		  - buf buffer.Appender - буфер для хранения закодированного значения.
	*/
	Encode(buf buffer.Appender)

	/*
		Command возвращает код команды, к которому относится тело сообщения.
	*/
	Command() command.Code

	/*
		IsValid возвращает валидность тела.
	*/
	IsValid() bool
}
