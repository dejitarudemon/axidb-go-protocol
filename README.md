[![CI](https://github.com/dejitarudemon/axidb-go-protocol/actions/workflows/ci.yml/badge.svg)](https://github.com/dejitarudemon/axidb-go-protocol/actions/workflows/ci.yml)
![coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/dejitarudemon/d4908af19df5c45057d630931ab6b8a7/raw/coverage.json)

# axidb-go-protocol

Go-библиотека кадров протокола [AxiDB](docs/specs.md): Hello (версия 0) и рабочая версия 1.

[Русский](#русский) · [English](#english)

---

## Русский

Кодирование, декодирование и проверка кадров AxiDB. Версия 0 согласовывает рабочую версию. Версия 1 — Handshake, Read, Write, Delete, Ping, Batch, ответы и ошибки.

### Спецификация

- [Русский текст](docs/specs.md)
- [English text](docs/specs.en.md)

### Установка

Требуется Go 1.27 или новее.

```bash
go get github.com/dejitarudemon/axidb-go-protocol@latest
```

### Пакеты

| Пакет | Назначение |
| --- | --- |
| `v0/builder`, `v0/decoder`, `v0/frame` | Hello-кадр: список поддерживаемых версий |
| `v0/body/bodies` | Тело Hello и пересечение версий (`Common`) |
| `v1/builder` | Сборка кадров Handshake, Read, Write, Delete, Ping, Batch и ответов |
| `v1/decoder` | Разбор кадра из `*bufio.Reader` |
| `v1/frame`, `v1/body/bodies` | Кадр и тела команд |
| `v1/value/values` | Типы значений: bytes, string, int, uint, float, JSON, массивы |
| `v1/compressor/compressors` | Сжатие тела: zstd, s2 |
| `v1/buffer` | Буфер кодирования (`buffer.Slice`) |
| `v1/err`, `v1/err/errs` | Локальные и протокольные ошибки |

### Туториал

Ниже — полный путь клиента: TCP, Hello, Handshake, запись, чтение, ping. Примеры собирают кадр в память; на сокете вместо `bytes.NewReader` используйте `bufio.NewReader(conn)`.

#### 1. Согласовать версию (Hello)

Клиент и сервер обмениваются кадрами версии 0. Общие версии — пересечение двух списков.

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

	// На сервере тот же кадр приходит из сети.
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
	fmt.Println("common:", mine.Common(peer)) // [1 2 3] в этом примере
}
```

`DecodePreamble` не потребляет байты: по версии `0` вызывайте `v0/decoder`, по версии `1` — `v1/decoder`.

Если пересечение пустое, клиент разрывает соединение.

#### 2. Рукопожатие (Handshake)

После Hello клиент отправляет кадр версии 1 с `RequestID = 0`. Сжатие в Handshake запрещено.

```go
package main

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/builder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func main() {
	fb := builder.NewFrameBuilder(1 << 20)
	var hash [32]byte // Argon2id от пароля; для анонимного входа — нули

	hs, err := fb.NewHandshake("user", hash, []fields.Compression{fields.Zstd, fields.S2})
	if err != nil {
		panic(err)
	}

	var out buffer.Slice
	out.Preallocate(hs.Size())
	if err := hs.Encode(&out, nil); err != nil {
		panic(err)
	}

	_ = out.Bytes() // отправить в соединение
}
```

Ответ сервера собирается через `NewHandshakeAnswer`. Дальше `RequestID` должен быть ненулевым.

#### 3. Запись и чтение

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

	fmt.Println(got.RequestID, got.Body) // 2 и тело Write

	read, err := fb.NewRead(1, fields.Key("code"))
	if err != nil {
		panic(err)
	}

	out.Clean()
	out.Preallocate(read.Size())
	if err := read.Encode(&out, nil); err != nil {
		panic(err)
	}

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

#### 4. Разобрать кадр из потока

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

#### 5. Сжатие тела

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

#### 6. Батч

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

### Примеры из спецификации

Готовые кадры из спецификации лежат в `v0/internal/specs` и `v1/internal/specs`. Тесты `TestFrame_Specs` и `TestDecoder_Specs` сверяют кодирование и разбор с этими байтами.

```text
Клиент → сервер, Hello:  0A DB 00 03 01 02 03 6E 38 99 00
Сервер → клиент, Hello:  0A DB 00 04 01 04 07 0B 3B 78 5D 98
Общие версии: 1
```

### Тесты

```bash
go test ./v0/... ./v1/...
```

---

## English

A Go library for [AxiDB](docs/specs.en.md) protocol frames: Hello (version 0) and working version 1.

### Specification

- [Russian](docs/specs.md)
- [English](docs/specs.en.md)

### Install

Go 1.27 or newer is required.

```bash
go get github.com/dejitarudemon/axidb-go-protocol@latest
```

### Packages

| Package | Role |
| --- | --- |
| `v0/builder`, `v0/decoder`, `v0/frame` | Hello frame: advertised protocol versions |
| `v0/body/bodies` | Hello body and version intersection (`Common`) |
| `v1/builder` | Handshake, Read, Write, Delete, Ping, Batch, and answers |
| `v1/decoder` | Parse one frame from a `*bufio.Reader` |
| `v1/frame`, `v1/body/bodies` | Frame and command bodies |
| `v1/value/values` | Value types: bytes, string, int, uint, float, JSON, arrays |
| `v1/compressor/compressors` | Body compression: zstd, s2 |
| `v1/buffer` | Encoding buffer (`buffer.Slice`) |
| `v1/err`, `v1/err/errs` | Local and protocol errors |

### Tutorial

A full client path: TCP, Hello, Handshake, write, read, ping. The snippets encode into memory; on a socket use `bufio.NewReader(conn)` instead of `bytes.NewReader`.

#### 1. Agree on a version (Hello)

Client and server exchange version-0 frames. Working versions are the intersection of the two lists.

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

	// On the server the same frame arrives from the network.
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
	fmt.Println("common:", mine.Common(peer)) // [1 2 3] in this example
}
```

`DecodePreamble` does not consume bytes: use `v0/decoder` for version `0` and `v1/decoder` for version `1`.

If the intersection is empty, the client closes the connection.

#### 2. Handshake

After Hello the client sends a version-1 frame with `RequestID = 0`. Compression is forbidden on Handshake.

```go
package main

import (
	"github.com/dejitarudemon/axidb-go-protocol/v1/builder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
)

func main() {
	fb := builder.NewFrameBuilder(1 << 20)
	var hash [32]byte // Argon2id of the password; zeros for anonymous login

	hs, err := fb.NewHandshake("user", hash, []fields.Compression{fields.Zstd, fields.S2})
	if err != nil {
		panic(err)
	}

	var out buffer.Slice
	out.Preallocate(hs.Size())
	if err := hs.Encode(&out, nil); err != nil {
		panic(err)
	}

	_ = out.Bytes() // write to the connection
}
```

Build the server reply with `NewHandshakeAnswer`. Later requests must use a non-zero `RequestID`.

#### 3. Write and read

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

	fmt.Println(got.RequestID, got.Body) // 2 and a Write body

	read, err := fb.NewRead(1, fields.Key("code"))
	if err != nil {
		panic(err)
	}

	out.Clean()
	out.Preallocate(read.Size())
	if err := read.Encode(&out, nil); err != nil {
		panic(err)
	}

	answer, err := fb.NewReadAnswer(1, values.String("IDDQD"))
	if err != nil {
		panic(err)
	}
	_ = answer
}
```

Other bodies:

```go
fb.NewDelete(3, fields.Key("another-key"))
fb.NewPing(4)
fb.NewWriteAnswer(2)
fb.NewDeleteAnswer(3)
fb.NewPingAnswer(4)
```

#### 4. Decode a frame from a stream

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

#### 5. Compress the body

Compression applies to the Body only. The checksum is computed over the compressed bytes. Handshake and Ping must not be compressed.

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

#### 6. Batch

A batch is not a transaction. Only Read, Write, and Delete are allowed. The builder assigns operation numbers starting at zero.

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

### Specification examples

Canonical frames from the spec live in `v0/internal/specs` and `v1/internal/specs`. `TestFrame_Specs` and `TestDecoder_Specs` check encode and decode against those bytes.

```text
Client → server, Hello:  0A DB 00 03 01 02 03 6E 38 99 00
Server → client, Hello:  0A DB 00 04 01 04 07 0B 3B 78 5D 98
Common versions: 1
```

### Tests

```bash
go test ./v0/... ./v1/...
```
