[![CI](https://github.com/dejitarudemon/axidb-go-protocol/actions/workflows/ci.yml/badge.svg)](https://github.com/dejitarudemon/axidb-go-protocol/actions/workflows/ci.yml)
![coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/dejitarudemon/d4908af19df5c45057d630931ab6b8a7/raw/coverage.json)

# axidb-go-protocol

A Go library for AxiDB protocol frames: Hello (version 0) and working version 1.

[English](#english) · [Русский](#русский)

---

## English

A Go library for AxiDB protocol frames: Hello (version 0) and working version 1.

After Hello the connection is not pinned to one version: each later frame may use any version both sides advertised.

### Docs

- [How it works](docs/how-it-works.en.md)
- [Tutorials](docs/tutorial.en.md)
- [Specification](docs/specs.en.md) · [Русский](docs/specs.md)

### Install

Go 1.27 or newer is required.

```bash
go get github.com/dejitarudemon/axidb-go-protocol@latest
```

### Example: TCP server

The server accepts Hello, replies with its versions, completes Handshake, then serves Ping, Read, Write, and Delete in memory.

```go
package main

import (
	"bufio"
	"log"
	"net"
	"sync"

	v0bodies "github.com/dejitarudemon/axidb-go-protocol/v0/body/bodies"
	v0builder "github.com/dejitarudemon/axidb-go-protocol/v0/builder"
	v0decoder "github.com/dejitarudemon/axidb-go-protocol/v0/decoder"
	v0fields "github.com/dejitarudemon/axidb-go-protocol/v0/fields"
	v0frame "github.com/dejitarudemon/axidb-go-protocol/v0/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/builder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/decoder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:4747")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listen", ln.Addr())

	for {
		c, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serve(c)
	}
}

type kv struct {
	mu   sync.Mutex
	data map[string]value.V
}

func serve(c net.Conn) {
	defer c.Close()

	r := bufio.NewReader(c)
	fb0 := v0builder.NewFrameBuilder(1 << 10)
	fb1 := builder.NewFrameBuilder(1 << 20)
	store := &kv{data: map[string]value.V{}}

	hello, err := v0decoder.NewDecoder().DecodeFrame(r)
	if err != nil {
		return
	}
	peer, ok := hello.Body.(v0bodies.Hello)
	if !ok {
		return
	}
	mine, err := fb0.NewHello([]v0fields.Version{1})
	if err != nil || len(mine.Body.(v0bodies.Hello).Common(peer)) == 0 {
		return
	}
	if err := writeV0(c, mine); err != nil {
		return
	}

	hs, err := decoder.NewDecoder(1<<20, nil).DecodeFrame(r)
	if err != nil {
		return
	}
	if _, ok := hs.Body.(bodies.Handshake); !ok {
		return
	}
	ans, err := fb1.NewHandshakeAnswer(nil)
	if err != nil {
		return
	}
	if err := writeV1(c, ans); err != nil {
		return
	}

	d1 := decoder.NewDecoder(1<<20, nil)
	for {
		f, err := d1.DecodeFrame(r)
		if err != nil {
			return
		}
		reply, err := store.handle(fb1, f)
		if err != nil {
			return
		}
		if err := writeV1(c, reply); err != nil {
			return
		}
	}
}

func (s *kv) handle(fb builder.FrameBuilder, f frame.Frame) (frame.Frame, error) {
	switch b := f.Body.(type) {
	case bodies.Ping:
		return fb.NewPingAnswer(f.RequestID)
	case bodies.Read:
		s.mu.Lock()
		v, ok := s.data[string(b)]
		s.mu.Unlock()
		if !ok {
			return fb.NewErrAnswer(f.RequestID, errs.NewErrorNotFound(fields.Key(b)))
		}
		return fb.NewReadAnswer(f.RequestID, v)
	case bodies.Write:
		s.mu.Lock()
		s.data[string(b.Key)] = b.Value
		s.mu.Unlock()
		return fb.NewWriteAnswer(f.RequestID)
	case bodies.Delete:
		s.mu.Lock()
		delete(s.data, string(b))
		s.mu.Unlock()
		return fb.NewDeleteAnswer(f.RequestID)
	default:
		cmd := fields.Command(0)
		if f.Body != nil {
			cmd = f.Body.Command()
		}
		return fb.NewErrAnswer(f.RequestID, errs.NewErrorUnsupportedCommand(cmd))
	}
}

func writeV0(c net.Conn, f v0frame.Frame) error {
	var buf buffer.Slice
	buf.Preallocate(f.Size())
	if err := f.Encode(&buf); err != nil {
		return err
	}
	_, err := c.Write(buf.Bytes())
	return err
}

func writeV1(c net.Conn, f frame.Frame) error {
	var buf buffer.Slice
	buf.Preallocate(f.Size())
	if err := f.Encode(&buf, nil); err != nil {
		return err
	}
	_, err := c.Write(buf.Bytes())
	return err
}
```

Client, compression, and batch are in the [tutorials](docs/tutorial.en.md).

### Tests and benchmarks

```bash
go test ./v0/... ./v1/...
go test -bench=. -benchmem ./v0/... ./v1/...
```

---

## Русский

Go-библиотека кадров протокола AxiDB: Hello (версия 0) и рабочая версия 1.

После Hello соединение не привязано к одной версии: каждый следующий кадр может использовать любую версию из пересечения списков.

### Документация

- [Как это работает](docs/how-it-works.md)
- [Туториалы](docs/tutorial.md)
- [Спецификация](docs/specs.md) · [English](docs/specs.en.md)

### Установка

Нужен Go 1.27 или новее.

```bash
go get github.com/dejitarudemon/axidb-go-protocol@latest
```

### Пример: TCP-сервер

Сервер принимает Hello, отвечает своими версиями, проходит Handshake и обслуживает Ping, Read, Write, Delete в памяти.

```go
package main

import (
	"bufio"
	"log"
	"net"
	"sync"

	v0bodies "github.com/dejitarudemon/axidb-go-protocol/v0/body/bodies"
	v0builder "github.com/dejitarudemon/axidb-go-protocol/v0/builder"
	v0decoder "github.com/dejitarudemon/axidb-go-protocol/v0/decoder"
	v0fields "github.com/dejitarudemon/axidb-go-protocol/v0/fields"
	v0frame "github.com/dejitarudemon/axidb-go-protocol/v0/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/builder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/decoder"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:4747")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listen", ln.Addr())

	for {
		c, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serve(c)
	}
}

type kv struct {
	mu   sync.Mutex
	data map[string]value.V
}

func serve(c net.Conn) {
	defer c.Close()

	r := bufio.NewReader(c)
	fb0 := v0builder.NewFrameBuilder(1 << 10)
	fb1 := builder.NewFrameBuilder(1 << 20)
	store := &kv{data: map[string]value.V{}}

	hello, err := v0decoder.NewDecoder().DecodeFrame(r)
	if err != nil {
		return
	}
	peer, ok := hello.Body.(v0bodies.Hello)
	if !ok {
		return
	}
	mine, err := fb0.NewHello([]v0fields.Version{1})
	if err != nil || len(mine.Body.(v0bodies.Hello).Common(peer)) == 0 {
		return
	}
	if err := writeV0(c, mine); err != nil {
		return
	}

	hs, err := decoder.NewDecoder(1<<20, nil).DecodeFrame(r)
	if err != nil {
		return
	}
	if _, ok := hs.Body.(bodies.Handshake); !ok {
		return
	}
	ans, err := fb1.NewHandshakeAnswer(nil)
	if err != nil {
		return
	}
	if err := writeV1(c, ans); err != nil {
		return
	}

	d1 := decoder.NewDecoder(1<<20, nil)
	for {
		f, err := d1.DecodeFrame(r)
		if err != nil {
			return
		}
		reply, err := store.handle(fb1, f)
		if err != nil {
			return
		}
		if err := writeV1(c, reply); err != nil {
			return
		}
	}
}

func (s *kv) handle(fb builder.FrameBuilder, f frame.Frame) (frame.Frame, error) {
	switch b := f.Body.(type) {
	case bodies.Ping:
		return fb.NewPingAnswer(f.RequestID)
	case bodies.Read:
		s.mu.Lock()
		v, ok := s.data[string(b)]
		s.mu.Unlock()
		if !ok {
			return fb.NewErrAnswer(f.RequestID, errs.NewErrorNotFound(fields.Key(b)))
		}
		return fb.NewReadAnswer(f.RequestID, v)
	case bodies.Write:
		s.mu.Lock()
		s.data[string(b.Key)] = b.Value
		s.mu.Unlock()
		return fb.NewWriteAnswer(f.RequestID)
	case bodies.Delete:
		s.mu.Lock()
		delete(s.data, string(b))
		s.mu.Unlock()
		return fb.NewDeleteAnswer(f.RequestID)
	default:
		cmd := fields.Command(0)
		if f.Body != nil {
			cmd = f.Body.Command()
		}
		return fb.NewErrAnswer(f.RequestID, errs.NewErrorUnsupportedCommand(cmd))
	}
}

func writeV0(c net.Conn, f v0frame.Frame) error {
	var buf buffer.Slice
	buf.Preallocate(f.Size())
	if err := f.Encode(&buf); err != nil {
		return err
	}
	_, err := c.Write(buf.Bytes())
	return err
}

func writeV1(c net.Conn, f frame.Frame) error {
	var buf buffer.Slice
	buf.Preallocate(f.Size())
	if err := f.Encode(&buf, nil); err != nil {
		return err
	}
	_, err := c.Write(buf.Bytes())
	return err
}
```

Клиент, сжатие и батч — в [туториалах](docs/tutorial.md).

### Тесты и бенчмарки

```bash
go test ./v0/... ./v1/...
go test -bench=. -benchmem ./v0/... ./v1/...
```

