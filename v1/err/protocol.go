package err

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/google/uuid"
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
	TracebackID() uuid.UUID

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
	Code() Code

	/*
		func IsValid возвращает валидность настоящей ошибки
		согласно спецификациям протокола v1.
	*/
	IsValid() error
}
