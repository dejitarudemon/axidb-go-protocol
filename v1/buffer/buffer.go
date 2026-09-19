/*
package buffer предназначен для объявления требований
к буферу, куда будут кодироваться кадры
*/
package buffer

/*
interface Appender предназначен для объявления требований
к буферу в частности добавления значений.
*/
type Appender interface {
	/*
		func Append предназначена для копирования []byte в буфер.

		Принимаемые значения:
			- v []byte - значение для копирования.
	*/
	Append(v []byte)

	/*
		func AppendUint8 предназначена для копирования uint8 в буфер.

		Принимаемые значения:
			- v uint8 - значение для копирования.
	*/
	AppendUint8(v uint8)

	/*
		func AppendUint16 предназначена для копирования uint16 в буфер.

		Принимаемые значения:
			- v uint16 - значение для копирования.
	*/
	AppendUint16(v uint16)

	/*
		func AppendUint32 предназначена для копирования uint32 в буфер.

		Принимаемые значения:
			- v uint32 - значение для копирования.
	*/
	AppendUint32(v uint32)

	/*
		func AppendUint64 предназначена для копирования uint64 в буфер.

		Принимаемые значения:
			- v uint64 - значение для копирования.
	*/
	AppendUint64(v uint64)

	/*
		func AppendString предназначена для копирования string в буфер.

		Принимаемые значения:
			- v string - значение для копирования.
	*/
	AppendString(v string)
}

/*
interface Buffer предназначен для расширения Appender в частности работы
с памятью.
*/
type Buffer interface {
	Appender

	/*
		func Preallocate предназначена для выделения буфера.

		Принимаемые значения:
			- size int - размер буфера.
	*/
	Preallocate(size int)

	/*
		func Clean предназначена для очищения буфера.
	*/
	Clean()

	/*
		func Bytes предназначена для получения значения буфера.

		Выходные значения:
			- []byte - значение в буфере.
	*/
	Bytes() []byte
}
