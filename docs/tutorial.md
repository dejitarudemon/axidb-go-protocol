Язык: [Русский](tutorial.md) · [English](tutorial.en.md)

# Туториалы

Пошаговые примеры кадров вне сети. Готовый TCP-сервер — в [README](../README.md#пример-tcp-сервер). Как устроен протокол — в [how-it-works](how-it-works.md).

Кадр в примерах собирается в память. На сокете вместо `bytes.NewReader` используйте `bufio.NewReader(conn)`.

## 1. Согласовать версию (Hello)

Клиент и сервер обмениваются кадрами версии 0. Общие версии — пересечение двух списков. `NewHello` отбрасывает версию 0, дубликаты и значения больше 255.

```go
package main

import (
	"bufio"
	"bytes"
	"fmt"

	v0bodies "github.com/dejitarudemon/axidb-go-protocol/v0/body/bodies"
	v0builder "github.com/dejitarudemon/axidb-go-protocol/v0/builder"
	v0decoder "github.com/dejitarudemon/axidb-go-protocol/v0/decoder"
	v0fields "github.com/dejitarudemon/axidb-go-protocol/v0/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
)

func main() {
	client := v0builder.NewFrameBuilder(1 << 10)
	hello, err := client.NewHello([]v0fields.Version{1, 2, 3})
	if err != nil {
		panic(err)
	}

	var out buffer.Slice
	out.Preallocate(hello.Size())
	if err := hello.Encode(&out); err != nil {
		panic(err)
	}

	d := v0decoder.NewDecoder()
	r := bufio.NewReader(bytes.NewReader(out.Bytes()))

	version, err := d.DecodePreamble(r)
	if err != nil {
		panic(err)
	}
	fmt.Println("preamble version:", version) // 0

	got, err := d.DecodeFrame(r)
	if err != nil {
		panic(err)
	}

	mine := hello.Body.(v0bodies.Hello)
	peer := got.Body.(v0bodies.Hello)
	fmt.Println("common:", mine.Common(peer))
}
```

`DecodePreamble` не потребляет байты: по версии `0` вызывайте `v0/decoder`, по версии `1` — `v1/decoder`. Если пересечение пустое, клиент разрывает соединение.

Пример из спецификации:

```text
Клиент → сервер:  0A DB 00 03 01 02 03 6E 38 99 00
Сервер → клиент:  0A DB 00 04 01 04 07 0B 3B 78 5D 98
Общие версии: 1
```

## 2. Рукопожатие (Handshake)

После Hello клиент отправляет кадр версии 1 с `RequestID = 0`. Сжатие в Handshake запрещено. Хеш — Argon2id от пароля; для анонимного входа — нули.

```go
package main

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/builder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func main() {
	fb := builder.NewFrameBuilder(1 << 20)
	var hash [32]byte

	hs, err := fb.NewHandshake("user", hash, []fields.Compression{fields.Zstd, fields.S2})
	if err != nil {
		panic(err)
	}

	var out buffer.Slice
	out.Preallocate(hs.Size())
	if err := hs.Encode(&out, nil); err != nil {
		panic(err)
	}

	_ = out.Bytes()
}
```

Ответ сервера — `NewHandshakeAnswer`. Дальше `RequestID` должен быть ненулевым.

## 3. Запись и чтение

```go
package main

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/dejitarudemon/axidb-go-protocol/v1/builder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/decoder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func main() {
	fb := builder.NewFrameBuilder(1 << 20)

	write, err := fb.NewWrite(2, fields.Key("code"), values.String("IDDQD"))
	if err != nil {
		panic(err)
	}

	var out buffer.Slice
	out.Preallocate(write.Size())
	if err := write.Encode(&out, nil); err != nil {
		panic(err)
	}

	d := decoder.NewDecoder(1<<20, nil)
	got, err := d.DecodeFrame(bufio.NewReader(bytes.NewReader(out.Bytes())))
	if err != nil {
		panic(err)
	}
	fmt.Println(got.RequestID, got.Body)

	read, err := fb.NewRead(1, fields.Key("code"))
	if err != nil {
		panic(err)
	}
	_ = read

	answer, err := fb.NewReadAnswer(1, values.String("IDDQD"))
	if err != nil {
		panic(err)
	}
	_ = answer
}
```

Другие тела:

```go
fb.NewDelete(3, fields.Key("another-key"))
fb.NewPing(4)
fb.NewWriteAnswer(2)
fb.NewDeleteAnswer(3)
fb.NewPingAnswer(4)
```

## 4. Разобрать кадр из потока

```go
package main

import (
	"bufio"
	"io"

	v0decoder "github.com/dejitarudemon/axidb-go-protocol/v0/decoder"
	v1decoder "github.com/dejitarudemon/axidb-go-protocol/v1/decoder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func readFrame(r *bufio.Reader) error {
	v0 := v0decoder.NewDecoder()
	version, err := v0.DecodePreamble(r)
	if err != nil {
		return err
	}

	if version == 0 {
		_, err = v0.DecodeFrame(r)
		return err
	}

	v1 := v1decoder.NewDecoder(fields.BodyLimit(1<<20), nil)
	_, err = v1.DecodeFrame(r)
	return err
}

func serve(r io.Reader) error {
	return readFrame(bufio.NewReader(r))
}
```

## 5. Сжатие тела

Сжатие применяется только к Body. Контрольная сумма считается по сжатым байтам. Handshake и Ping сжимать нельзя.

```go
package main

import (
	"bufio"
	"bytes"

	"github.com/dejitarudemon/axidb-go-protocol/v1/builder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor"
	"github.com/dejitarudemon/axidb-go-protocol/v1/compressor/compressors"
	"github.com/dejitarudemon/axidb-go-protocol/v1/decoder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func main() {
	zstd, err := compressors.NewZstd(1 << 20)
	if err != nil {
		panic(err)
	}

	fb := builder.NewFrameBuilder(1 << 20)
	f, err := fb.NewWrite(2, fields.Key("blob"), values.Bytes(make([]byte, 256)))
	if err != nil {
		panic(err)
	}

	var out buffer.Slice
	if err := f.Encode(&out, zstd); err != nil {
		panic(err)
	}

	d := decoder.NewDecoder(1<<20, []compressor.Compressor{zstd})
	_, err = d.DecodeFrame(bufio.NewReader(bytes.NewReader(out.Bytes())))
	if err != nil {
		panic(err)
	}
}
```

## 6. Батч

Батч — не транзакция. Допустимы только Read, Write и Delete. Номера операций назначает builder, начиная с нуля.

```go
package main

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/builder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func main() {
	fb := builder.NewFrameBuilder(1 << 20)
	batch := builder.NewBatchRequestsBuilder().
		SequentialExecution(true).
		InterruptAfterError(false).
		OneAnswer(true).
		AddRead(fields.Key("key")).
		AddWrite(fields.Key("another-key"), values.String("data"))

	f, err := fb.NewBatch(1, *batch)
	if err != nil {
		panic(err)
	}
	_ = f
}
```

Флаги:

| Флаг | 1 | 0 |
| --- | --- | --- |
| Sequential | по номеру | параллельно |
| Interrupt after error | остальные запросы получают Request Interrupted | остальные продолжаются |
| One answer | один общий ответ | ответы по мере готовности |

Готовые кадры из спецификации — в `v0/internal/specs` и `v1/internal/specs`.
