package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
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

func readLineWithPrefix(r *bufio.Reader, prefixes ...string) (prefix, str string, err error) {
	var line []byte
	if line, _, err = r.ReadLine(); err != nil {
		return
	}
	str = strings.TrimSpace(string(line))
	for _, pp := range prefixes {
		if strings.HasPrefix(str, pp) {
			return pp, str[len(pp):], nil
		}
	}
	return "", "", fmt.Errorf("Line expected to start with one of the prefixes: %v, the actual line is: '%s'", prefixes, str)
}

func consumeLine(r *bufio.Reader, want string) (err error) {
	var str string
	if _, str, err = readLineWithPrefix(r, want); err != nil {
		return
	}
	if str != "" {
		return fmt.Errorf("Line contains unexpected symbols after the right prefix: '%s', symbols: '%s'", want, str)
	}
	return nil
}

func readAsciiSTL(data []byte) (res []STLTriangle, err error) {
	r := bufio.NewReader(bytes.NewBuffer(data))
	if err = consumeLine(r, "solid object"); err != nil {
		return
	}
	for {
		var prefix, str string
		var t STLTriangle
		if prefix, str, err = readLineWithPrefix(r, "facet normal ", "endsolid object"); err != nil {
			if err == io.EOF {
				return res, nil
			}
			return nil, err
		}
		if prefix == "endsolid object" {
			return
		}
		fields := strings.Fields(str)
		if len(fields) != 3 {
			return nil, fmt.Errorf("Normal definition is broken: '%s'", str)
		}
		for i := 0; i < 3; i++ {
			var v float64
			if v, err = strconv.ParseFloat(fields[i], 32); err != nil {
				return nil, err
			}
			t.n[i] = float32(v)
		}
		if err = consumeLine(r, "outer loop"); err != nil {
			return nil, err
		}
		for i := 0; i < 3; i++ {
			if _, str, err = readLineWithPrefix(r, "vertex "); err != nil {
				return nil, err
			}
			fields = strings.Fields(str)
			if len(fields) != 3 {
				return nil, fmt.Errorf("Vertex definition is broken: '%s'", str)
			}
			for j := 0; j < 3; j++ {
				var v float64
				if v, err = strconv.ParseFloat(fields[j], 32); err != nil {
					return nil, err
				}
				t.v[i][j] = float32(v)
			}
		}
		if err = consumeLine(r, "endloop"); err != nil {
			return nil, err
		}
		if err = consumeLine(r, "endfacet"); err != nil {
			return nil, err
		}

		res = append(res, t)

	}
	return
}

func ReadSTL(r io.Reader) (t []STLTriangle, err error) {
	var data []byte
	if data, err = ioutil.ReadAll(r); err != nil {
		return
	}
	if len(data) < 5 {
		return nil, fmt.Errorf("The file is too short: %d bytes", len(data))
	}
	magic := data[:5]
	if string(magic) == "solid" {
		return readAsciiSTL(data)
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
