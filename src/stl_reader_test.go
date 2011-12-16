package main

import (
	"os"
	"testing"
)

type readSTLTest struct {
	filename string
	count    int
}

var readSTLTests = []readSTLTest{
	{"data/cylinder.bin.stl", 326},
	// TODO(krasin): implement ASCII STL reader
	//	{"data/cylinder.stl", 326},
}

func TestReadSTL(t *testing.T) {
	for _, test := range readSTLTests {
		f, err := os.Open(test.filename)
		if err != nil {
			t.Fatalf("os.Open(\"%v\"): %v", test.filename, err)
		}
		defer f.Close()
		stl, err := ReadSTL(f)
		if err != nil {
			t.Fatalf("ReadSTL: %v", err)
		}
		if len(stl) != test.count {
			t.Errorf("Wrong number of triangles. Expected: %d, got: %d", test.count, len(stl))
		}
	}
}
