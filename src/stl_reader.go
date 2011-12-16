package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unsafe"
)

type STLPoint [3]float32

type STLTriangle struct {
	n STLPoint
	v [3]STLPoint
}

func readSTLPoint(a []byte, p *STLPoint) []byte {
	for i := 0; i < 3; i++ {
		cur := uint32(a[0]) + uint32(a[1])<<8 + uint32(a[2])<<16 + uint32(a[3])<<24
		p[i] = *(*float32)(unsafe.Pointer(&cur))
		a = a[4:]
	}
	return a
}

func ReadSTL(r io.Reader) (t []STLTriangle, err os.Error) {
	var data []byte
	if data, err = ioutil.ReadAll(r); err != nil {
		return
	}
	magic := data[:5]
	if string(magic) == "solid" {
		panic("ReadSTL: ascii format is not implemented")
	}
	// Skip STL header
	data = data[80:]
	n := uint32(data[0]) + uint32(data[1])<<8 + uint32(data[2])<<16 + uint32(data[3])<<24
	data = data[4:]

	if len(data) < int(SizeOfSTLTriangle*n) {
		return nil, fmt.Errorf("ReadSTL: unexpected end of file: want %d bytes to read triangle data, but only %d bytes is available", SizeOfSTLTriangle*n, len(data))
	}
	for i := 0; i < int(n); i++ {
		var cur STLTriangle
		data = readSTLPoint(data, &cur.n)
		for j := 0; j < 3; j++ {
			data = readSTLPoint(data, &cur.v[j])
		}
		data = data[2:]
		t = append(t, cur)
	}
	return
}
