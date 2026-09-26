// Package specs holds the example Hello frames from docs/specs.md shared by protocol v0 tests.
package specs

import (
	"github.com/dejitarudemon/axidb-go-protocol/v0/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v0/fields"
	"github.com/dejitarudemon/axidb-go-protocol/v0/frame"
)

// Frame is an example frame together with its exact wire encoding.
type Frame struct {
	Name    string
	Encoded []byte
	Frame   frame.Frame
}

// Frames lists every Hello example from the specification.
var Frames = []Frame{
	{
		Name: "client hello",
		Encoded: []byte{
			0x0A, 0xDB, 0x00, 0x03, 0x01, 0x02, 0x03, 0x6E, 0x38, 0x99, 0x00,
		},
		Frame: frame.Frame{
			Body: bodies.Hello{
				Versions: []fields.Version{1, 2, 3},
			},
		},
	},
	{
		Name: "server hello",
		Encoded: []byte{
			0x0A, 0xDB, 0x00, 0x04, 0x01, 0x04, 0x07, 0x0B, 0x3B, 0x78, 0x5D, 0x98,
		},
		Frame: frame.Frame{
			Body: bodies.Hello{
				Versions: []fields.Version{1, 4, 7, 11},
			},
		},
	},
}
