package err

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

/*
TracebackIDFieldSize представляет размер в байтах,
отведенный для хранения Traceback ID ошибки в сообщении.
*/
const (
	TracebackIDFieldSize = 16
)

/*
interface Error предназначен для представления
ошибки протокола согласно спецификации протокола v1.
*/
type ProtocolError interface {
	error

	/*
		func Size возвращает размер сообщения ошибки в байтах.
	*/
	Size() int

	/*
		func TracebackID возвращает Traceback ID ошибки.
	*/
	TracebackID() fields.TracebackID

	/*
		func Encode предназначена для кодирования сообщения об ошибке
		в сообщении.

		Принимаемые параметры:
		  - buf buffer.Appender - буфер для хранения закодированного значения.
	*/
	Encode(buf buffer.Appender)

	/*
		func Code возвращает код ошибки.
	*/
	Code() fields.Error
}
