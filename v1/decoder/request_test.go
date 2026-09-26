package decoder

import (
	"errors"
	"math"
	"runtime"
	"testing"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body"
	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/err/errs"
	"github.com/dejitarudemon/axidb-go-protocol/v1/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func batchRequestBytes(number uint32, command fields.Command, nested []byte) []byte {
	return cat(u32(number), []byte{byte(command)}, u32(uint32(len(nested))), nested)
}

func TestDecoder_decodeBody_Requests(t *testing.T) {
	tests := []struct {
		name string
		want body.Body
	}{
		{"ping", bodies.Ping{}},
		{"read", bodies.Read("key")},
		{"read empty key", bodies.Read("")},
		{"delete", bodies.Delete("key")},
		{"handshake empty", bodies.NewHandshake("", [32]byte{}, nil)},
		{"handshake", bodies.NewHandshake("user", testHash, []fields.Compression{fields.Zstd, fields.S2})},
		{"write empty key", bodies.Write{Key: fields.Key{}, Value: values.Uint(1)}},
		{"write json", bodies.Write{Key: fields.Key("k"), Value: values.JSON(`{"a":1}`)}},
		{"write float", bodies.Write{Key: fields.Key("k"), Value: values.Float(0.5)}},
		{"batch empty", bodies.Batch{Requests: []bodies.Request{}}},
		{"batch", bodies.Batch{
			IsSequentialExecution: true,
			InterruptAfterError:   true,
			IsOneAnswer:           true,
			Requests: []bodies.Request{
				{Number: 0, Body: bodies.Read("a")},
				{Number: 1, Body: bodies.Write{Key: fields.Key("b"), Value: values.Int(-2)}},
				{Number: 2, Body: bodies.Delete("c")},
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCursor(encodeBody(t, tt.want))
			got, e := NewDecoder(1024, nil).decodeBody(c, tt.want.Command())
			assertDecoded(t, got, e, c, tt.want)
		})
	}
}

func TestDecoder_decodeBody_RequestErrs(t *testing.T) {
	hash := make([]byte, bodies.HashFieldSize)
	readReq := batchRequestBytes(0, fields.Read, []byte("k"))

	tests := []struct {
		name    string
		command fields.Command
		data    []byte
	}{
		{"handshake empty", fields.Handshake, nil},
		{"handshake login shorter than len", fields.Handshake, cat(u32(5), []byte("user"))},
		{"handshake without hash", fields.Handshake, cat(u32(0), hash[:31])},
		{"handshake without compressions len", fields.Handshake, cat(u32(0), hash)},
		{"handshake compressions shorter than len", fields.Handshake, cat(u32(0), hash, []byte{2, 1})},
		{"write empty", fields.Write, nil},
		{"write key shorter than len", fields.Write, cat(u32(4), []byte("key"))},
		{"write without type", fields.Write, cat(u32(1), []byte("k"))},
		{"write unknown type", fields.Write, cat(u32(1), []byte("k"), []byte{0xFF})},
		{"write short value", fields.Write, cat(u32(1), []byte("k"), []byte{byte(fields.Int)}, u32(0))},
		{"batch empty", fields.Batch, nil},
		{"batch without len", fields.Batch, []byte{0x00, 0x00}},
		{"batch missing request", fields.Batch, cat([]byte{0x00}, u32(2), readReq)},
		{"batch huge len", fields.Batch, cat([]byte{0x00}, u32(math.MaxUint32), readReq)},
		{"batch request without number", fields.Batch, cat([]byte{0x00}, u32(1), []byte{0x00})},
		{"batch request without command", fields.Batch, cat([]byte{0x00}, u32(1), u32(0))},
		{"batch request without body len", fields.Batch, cat([]byte{0x00}, u32(1), u32(0), []byte{byte(fields.Read)})},
		{"batch request body longer than batch", fields.Batch, cat([]byte{0x00}, u32(1), u32(0), []byte{byte(fields.Read)}, u32(100), []byte("k"))},
		{"batch request body len overflows", fields.Batch, cat([]byte{0x00}, u32(1), u32(0), []byte{byte(fields.Read)}, u32(math.MaxUint32))},
		{"batch request malformed body", fields.Batch, cat([]byte{0x00}, u32(1), batchRequestBytes(0, fields.Write, u32(1)))},
		{"batch request trailing bytes", fields.Batch, cat([]byte{0x00}, u32(1), batchRequestBytes(0, fields.Ping, []byte{0x01}))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, e := NewDecoder(1024, nil).decodeBody(newCursor(tt.data), tt.command)
			assertMalformed(t, e)
		})
	}
}

func TestDecoder_batch_UnexpectedCommand(t *testing.T) {
	tests := []struct {
		name    string
		command fields.Command
		want    any
	}{
		{"batch in batch", fields.Batch, errs.ErrorUnexpectedCommandInBatch{}},
		{"unsupported command in batch", fields.Command(0xFF), errs.ErrorUnsupportedCommand{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := cat([]byte{0x00}, u32(1), batchRequestBytes(3, tt.command, nil))
			_, e := NewDecoder(1024, nil).decodeBody(newCursor(data), fields.Batch)
			if !errors.As(e, &tt.want) {
				t.Fatalf("got %v, want %T", e, tt.want)
			}
		})
	}
}

func allocatedBytes(f func()) uint64 {
	var before, after runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&before)
	f()
	runtime.ReadMemStats(&after)

	return after.TotalAlloc - before.TotalAlloc
}

func TestDecoder_DoesNotTrustDeclaredCounts(t *testing.T) {
	tests := []struct {
		name    string
		command fields.Command
		data    []byte
	}{
		{"batch", fields.Batch, cat([]byte{0x00}, u32(math.MaxUint32), batchRequestBytes(0, fields.Ping, nil))},
		{"batch answer", fields.Answer, cat([]byte{resultOK, byte(fields.Batch)}, u32(math.MaxUint32))},
		{"typed array", fields.Write, cat(u32(0), []byte{byte(fields.TypedArray)}, u32(math.MaxUint32), []byte{byte(fields.Int)})},
		{"untyped array", fields.Write, cat(u32(0), []byte{byte(fields.UntypedArray)}, u32(math.MaxUint32))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDecoder(1024, nil)

			allocated := allocatedBytes(func() {
				_, e := d.decodeBody(newCursor(tt.data), tt.command)
				assertMalformed(t, e)
			})

			if allocated > 1<<20 {
				t.Errorf("allocated %v bytes for a %v-byte body", allocated, len(tt.data))
			}
		})
	}
}
