package main

import (
	"fmt"
	"uuid"

	"github.com/dejitarudemon/axidb-go-protocol/v1/body/bodies"
	"github.com/dejitarudemon/axidb-go-protocol/v1/buffer"
	"github.com/dejitarudemon/axidb-go-protocol/v1/frame"
	"github.com/dejitarudemon/axidb-go-protocol/v1/value/values"
)

func main() {
	f := frame.Frame{
		RequestID: 1,
		Body: bodies.ReadAnswer{
			Value: values.String("some-data"),
		},
	}

	if err := f.IsValid(); err != nil {
		fmt.Println(err)
		return
	}

	buf := buffer.Slice{}
	buf.Preallocate(f.Size())

	if err := f.Encode(&buf, nil); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("% X\n", buf.Bytes())

	u := uuid.MustParse("6bae3680-7695-4742-92af-0a6fec339825")
	for _, char := range u {
		fmt.Printf("0x%X, ", char)
	}

	fmt.Println()
}
