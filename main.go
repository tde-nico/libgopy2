package main

import (
	"lib/libgopy"
	"log"
)

func main() {
	libgopy.Init()
	defer libgopy.Finalize()

	err := libgopy.Load("libtests.test_script2")
	if err != nil {
		log.Printf("Error: %v", err)
	}

	ret, err := libgopy.Call("func9",
		float64(71.5), float32(3.14),
		int64(3), int32(4), int16(5), int8(6), int(10),
		uint64(3), uint32(4), uint16(5), uint8(6), uint(10),
		rune(65),
		byte('A'),
		[]uint8("World"),
		[]byte("World"),
		"Hello",
		"Hello\x00World",
		[]byte("Hello\x00World"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		log.Printf("Return |%T| -> |%+v|", ret, ret)
	}
}
