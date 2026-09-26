package buffer

import (
	"bytes"
	"math"
	"testing"
)

func TestSlice_Preallocate(t *testing.T) {
	tests := []struct {
		name string
		size int
		want int
	}{
		{"zero", 0, 0},
		{"one byte", 1, 1},
		{"one kilobyte", 1 << 10, 1 << 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Slice{}
			s.Preallocate(tt.size)

			if got := cap(s.data); got != tt.want {
				t.Errorf("cap = %v, want %v", got, tt.want)
			}

			if got := len(s.data); got != 0 {
				t.Errorf("len = %v, want 0", got)
			}
		})
	}
}

func TestSlice_PreallocateAfterPreallocate(t *testing.T) {
	tests := []struct {
		name   string
		first  int
		second int
		want   int
	}{
		{"grow from zero", 0, 1, 1},
		{"shrink to zero", 1, 0, 1},
		{"grow", 1 << 10, 1 << 11, 1 << 11},
		{"shrink", 1 << 11, 1 << 10, 1 << 11},
		{"same size", 1, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Slice{}
			s.Preallocate(tt.first)
			s.Preallocate(tt.second)

			if got := cap(s.data); got != tt.want {
				t.Errorf("cap = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlice_PreallocateAfterAppend(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		size    int
		want    []byte
		wantCap int
	}{
		{"fits in capacity", []byte{0x01, 0x02}, 1, []byte{0x01, 0x02}, 2},
		{"exceeds capacity", []byte{0x01, 0x02}, 1 << 10, []byte{}, 1 << 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Slice{}
			s.Preallocate(len(tt.data))
			s.Append(tt.data)
			s.Preallocate(tt.size)

			if got := s.Bytes(); !bytes.Equal(got, tt.want) {
				t.Errorf("Bytes() = % X, want % X", got, tt.want)
			}

			if got := cap(s.data); got != tt.wantCap {
				t.Errorf("cap = %v, want %v", got, tt.wantCap)
			}
		})
	}
}

func TestSlice_PreallocateOvercrowded(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want []byte
	}{
		{"one byte", []byte{0x01}, []byte{0x01, 0x01}},
		{"one kilobyte", bytes.Repeat([]byte{0xFF}, 1<<10), bytes.Repeat([]byte{0xFF}, 1<<11)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Slice{}
			s.Preallocate(len(tt.data))
			s.Append(tt.data)
			s.Append(tt.data)

			if got := s.Bytes(); !bytes.Equal(got, tt.want) {
				t.Errorf("Bytes() = % X, want % X", got, tt.want)
			}
		})
	}
}

func TestSlice_Clean(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"empty", []byte{}},
		{"one byte", []byte{0x01}},
		{"one kilobyte", bytes.Repeat([]byte{0xFF}, 1<<10)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Slice{data: bytes.Clone(tt.data)}
			wantCap := cap(s.data)
			s.Clean()

			if got := s.Bytes(); len(got) != 0 {
				t.Errorf("Bytes() = % X, want empty", got)
			}

			if got := cap(s.data); got != wantCap {
				t.Errorf("cap = %v, want %v", got, wantCap)
			}

			s.AppendUint8(0xAB)
			if got, want := s.Bytes(), []byte{0xAB}; !bytes.Equal(got, want) {
				t.Errorf("Bytes() after append = % X, want % X", got, want)
			}
		})
	}
}

func TestSlice_Append(t *testing.T) {
	tests := []struct {
		name   string
		append func(*Slice)
		want   []byte
	}{
		{"nothing", func(*Slice) {}, []byte{}},
		{"nil bytes", func(s *Slice) { s.Append(nil) }, []byte{}},
		{"empty bytes", func(s *Slice) { s.Append([]byte{}) }, []byte{}},
		{"one byte", func(s *Slice) { s.Append([]byte{0x01}) }, []byte{0x01}},
		{"one kilobyte", func(s *Slice) { s.Append(bytes.Repeat([]byte{0xFF}, 1<<10)) }, bytes.Repeat([]byte{0xFF}, 1<<10)},
		{"empty string", func(s *Slice) { s.AppendString("") }, []byte{}},
		{"ascii string", func(s *Slice) { s.AppendString("hello world") }, []byte{0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64}},
		{"utf-8 string", func(s *Slice) { s.AppendString("привет мир") }, []byte{0xD0, 0xBF, 0xD1, 0x80, 0xD0, 0xB8, 0xD0, 0xB2, 0xD0, 0xB5, 0xD1, 0x82, 0x20, 0xD0, 0xBC, 0xD0, 0xB8, 0xD1, 0x80}},
		{"uint8 zero", func(s *Slice) { s.AppendUint8(0) }, []byte{0x00}},
		{"uint8 one", func(s *Slice) { s.AppendUint8(1) }, []byte{0x01}},
		{"uint8 max", func(s *Slice) { s.AppendUint8(math.MaxUint8) }, []byte{0xFF}},
		{"uint16 zero", func(s *Slice) { s.AppendUint16(0) }, []byte{0x00, 0x00}},
		{"uint16 one", func(s *Slice) { s.AppendUint16(1) }, []byte{0x00, 0x01}},
		{"uint16 big-endian", func(s *Slice) { s.AppendUint16(0x0102) }, []byte{0x01, 0x02}},
		{"uint16 max", func(s *Slice) { s.AppendUint16(math.MaxUint16) }, []byte{0xFF, 0xFF}},
		{"uint32 zero", func(s *Slice) { s.AppendUint32(0) }, []byte{0x00, 0x00, 0x00, 0x00}},
		{"uint32 one", func(s *Slice) { s.AppendUint32(1) }, []byte{0x00, 0x00, 0x00, 0x01}},
		{"uint32 big-endian", func(s *Slice) { s.AppendUint32(0x01020304) }, []byte{0x01, 0x02, 0x03, 0x04}},
		{"uint32 max", func(s *Slice) { s.AppendUint32(math.MaxUint32) }, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
		{"uint64 zero", func(s *Slice) { s.AppendUint64(0) }, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"uint64 one", func(s *Slice) { s.AppendUint64(1) }, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}},
		{"uint64 big-endian", func(s *Slice) { s.AppendUint64(0x0102030405060708) }, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}},
		{"uint64 max", func(s *Slice) { s.AppendUint64(math.MaxUint64) }, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{
			"mixed sequence",
			func(s *Slice) {
				s.AppendUint8(0x01)
				s.AppendUint16(0x0203)
				s.AppendUint32(0x04050607)
				s.AppendUint64(0x08090A0B0C0D0E0F)
				s.AppendString("ab")
				s.Append([]byte{0x10, 0x11})
			},
			[]byte{
				0x01,
				0x02, 0x03,
				0x04, 0x05, 0x06, 0x07,
				0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
				0x61, 0x62,
				0x10, 0x11,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allocations := []struct {
				name string
				size int
			}{
				{"without preallocate", -1},
				{"with preallocate", len(tt.want)},
			}

			for _, a := range allocations {
				t.Run(a.name, func(t *testing.T) {
					s := Slice{}
					if a.size >= 0 {
						s.Preallocate(a.size)
					}

					tt.append(&s)

					if got := s.Bytes(); !bytes.Equal(got, tt.want) {
						t.Errorf("Bytes() = % X, want % X", got, tt.want)
					}
				})
			}
		})
	}
}

func TestSlice_AppendCopiesInput(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"one byte", []byte{0x01}},
		{"one kilobyte", bytes.Repeat([]byte{0xFF}, 1<<10)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := bytes.Clone(tt.data)
			s := Slice{}
			s.Append(in)
			in[0] ^= 0xFF

			if got := s.Bytes(); !bytes.Equal(got, tt.data) {
				t.Errorf("Bytes() = % X, want % X", got, tt.data)
			}
		})
	}
}

func TestSlice_BytesReturnsCopy(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		preallocate int
	}{
		{"one byte", []byte{0x01}, 0},
		{"one kilobyte", bytes.Repeat([]byte{0xFF}, 1<<10), 0},
		{"spare capacity", []byte{0x01, 0x02}, 1 << 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Slice{}
			s.Preallocate(tt.preallocate)
			s.Append(tt.data)

			got := s.Bytes()
			got[0] ^= 0xFF

			if again := s.Bytes(); !bytes.Equal(again, tt.data) {
				t.Errorf("Bytes() after mutating previous result = % X, want % X", again, tt.data)
			}
		})
	}
}
