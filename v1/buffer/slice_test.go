package buffer

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func TestSlice_Preallocate(t *testing.T) {
	tests := []*struct {
		want int
	}{
		{0},
		{1},
		{1 << 10},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSlice_Preallocate_%v", i),
			func(t *testing.T) {
				slice := Slice{}
				slice.Preallocate(tt.want)

				if cap(slice.data) != tt.want {
					t.Errorf("preallocate failure: got %v cap want %v", cap(slice.data), tt.want)
				}
			},
		)
	}
}

func TestSlice_PreallocateAfterPreallocate(t *testing.T) {
	tests := []*struct {
		first  int
		second int
		want   int
	}{
		{0, 1, 1},
		{1, 0, 1},
		{1 << 10, 1 << 11, 1 << 11},
		{1 << 11, 1 << 10, 1 << 11},
		{1, 1, 1},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSlice_PreallocateAfterPreallocate_%v", i),
			func(t *testing.T) {
				slice := Slice{}
				slice.Preallocate(tt.first)
				slice.Preallocate(tt.second)

				if cap(slice.data) != tt.want {
					t.Errorf("preallocate failure: got %v cap want %v", cap(slice.data), tt.want)
				}
			},
		)
	}
}

func TestSlice_Clean(t *testing.T) {
	tests := []struct {
		s *Slice
	}{
		{&Slice{data: []byte{0x01}}},
		{&Slice{data: bytes.Repeat([]byte{0xFF}, 1024)}},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSlice_Clean_%v", i),
			func(t *testing.T) {
				tt.s.Clean()

				if cap(tt.s.data) == 0 {
					t.Errorf("clean failure: cap is 0")
				}

				if len(tt.s.data) != 0 {
					t.Errorf("clean failure: len is %v", len(tt.s.data))
				}
			},
		)
	}
}

func TestSlice_Append(t *testing.T) {
	tests := []struct {
		want     []byte
		allocate bool
	}{
		{[]byte{}, false},
		{[]byte{0x01}, false},
		{bytes.Repeat([]byte{0xFF}, 1024), false},

		{[]byte{}, true},
		{[]byte{0x01}, true},
		{bytes.Repeat([]byte{0xFF}, 1024), true},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSlice_Append_%v", i),
			func(t *testing.T) {
				slice := Slice{}

				if tt.allocate {
					slice.Preallocate(len(tt.want))
				}

				slice.Append(tt.want)

				if len(slice.data) != len(tt.want) {
					t.Errorf("append failure: len is %v, want %v", len(slice.data), len(tt.want))
				}
			},
		)
	}
}

func TestSlice_AppendString(t *testing.T) {
	tests := []struct {
		want     string
		allocate bool
	}{
		{"", false},
		{"a", false},
		{"hello world", false},
		{"хуйлоу ворлд", false},

		{"", true},
		{"a", true},
		{"hello world", true},
		{"хуйлоу ворлд", true},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSlice_AppendString_%v", i),
			func(t *testing.T) {
				slice := Slice{}

				if tt.allocate {
					slice.Preallocate(len(tt.want))
				}

				slice.AppendString(tt.want)

				if len(slice.data) != len(tt.want) {
					t.Errorf("append string failure: len is %v, want %v", len(slice.data), len(tt.want))
				}
			},
		)
	}
}

func TestSlice_AppendUint8(t *testing.T) {
	tests := []struct {
		want     uint8
		allocate bool
	}{
		{0, false},
		{1, false},
		{255, false},

		{0, true},
		{1, true},
		{255, true},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSlice_AppendUint8_%v", i),
			func(t *testing.T) {
				slice := Slice{}

				if tt.allocate {
					slice.Preallocate(1)
				}

				slice.AppendUint8(tt.want)

				if len(slice.data) != 1 {
					t.Errorf("append uint8 failure: len is %v, want 1", len(slice.data))
				}
			},
		)
	}
}

func TestSlice_AppendUint16(t *testing.T) {
	tests := []struct {
		want     uint16
		allocate bool
	}{
		{0, false},
		{1, false},
		{math.MaxUint16, false},

		{0, true},
		{1, true},
		{math.MaxUint16, true},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSlice_AppendUint16_%v", i),
			func(t *testing.T) {
				slice := Slice{}

				if tt.allocate {
					slice.Preallocate(2)
				}

				slice.AppendUint16(tt.want)

				if len(slice.data) != 2 {
					t.Errorf("append uint16 failure: len is %v, want 2", len(slice.data))
				}
			},
		)
	}
}

func TestSlice_AppendUint32(t *testing.T) {
	tests := []struct {
		want     uint32
		allocate bool
	}{
		{0, false},
		{1, false},
		{math.MaxUint32, false},

		{0, true},
		{1, true},
		{math.MaxUint32, true},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSlice_AppendUint32_%v", i),
			func(t *testing.T) {
				slice := Slice{}

				if tt.allocate {
					slice.Preallocate(4)
				}

				slice.AppendUint32(tt.want)

				if len(slice.data) != 4 {
					t.Errorf("append uint32 failure: len is %v, want 4", len(slice.data))
				}
			},
		)
	}
}

func TestSlice_AppendUint64(t *testing.T) {
	tests := []struct {
		want     uint64
		allocate bool
	}{
		{0, false},
		{1, false},
		{math.MaxUint64, false},

		{0, true},
		{1, true},
		{math.MaxUint64, true},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSlice_AppendUint64_%v", i),
			func(t *testing.T) {
				slice := Slice{}

				if tt.allocate {
					slice.Preallocate(8)
				}

				slice.AppendUint64(tt.want)

				if len(slice.data) != 8 {
					t.Errorf("append uint32 failure: len is %v, want 8", len(slice.data))
				}
			},
		)
	}
}

func TestSlice_Bytes(t *testing.T) {
	tests := []struct {
		want     []byte
		allocate bool
	}{
		{[]byte{}, false},
		{[]byte{0x01}, false},
		{bytes.Repeat([]byte{0xFF}, 1024), false},

		{[]byte{}, true},
		{[]byte{0x01}, true},
		{bytes.Repeat([]byte{0xFF}, 1024), true},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSlice_Bytes_%v", i),
			func(t *testing.T) {
				slice := Slice{}

				if tt.allocate {
					slice.Preallocate(len(tt.want))
				}

				slice.Append(tt.want)

				if len(slice.Bytes()) != len(tt.want) {
					t.Errorf("bytes() failure: len is %v, want %v", len(slice.Bytes()), len(tt.want))
				}
			},
		)
	}
}

func TestSlice_PreallocateOvercrowded(t *testing.T) {
	tests := []struct {
		want []byte
	}{
		{[]byte{0x01}},
		{bytes.Repeat([]byte{0xFF}, 1024)},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("TestSlice_PreallocateOvercrowded_%v", i),
			func(t *testing.T) {
				slice := Slice{}
				slice.Preallocate(len(tt.want))

				slice.Append(tt.want)
				slice.Append(tt.want)

				if len(slice.Bytes()) != len(tt.want)*2 {
					t.Errorf("preallocate overcrowded failure: len is %v, want %v", len(slice.Bytes()), len(tt.want)*2)
				}
			},
		)
	}
}
