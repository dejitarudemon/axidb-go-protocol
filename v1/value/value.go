package value

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

/*
interface V предназначен для объявления требований
к данным (Value).
*/
type V interface {
	/*
		func Size возвращает размер данных в байтах.
	*/
	Size() int

	/*
		func IsValid предназначена для проверки данных.
		Проверки: зависят от конкретной реализации интерфейса.

		Возвращаемые ошибки:
		 1. ErrorMalformedValue.
	*/
	IsValid() error

	/*
		func Encode предназначена для кодирования данных
		в сообщении.

		Принимааемые параметры:
		  - buf buffer.Appender - буфер для хранения закодированного значения.
	*/
	Encode(buf buffer.Appender)

	/*
		func Type возвращает код типа данных.
	*/
	Type() fields.Type
}
