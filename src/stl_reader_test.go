package main

import (
	"math"
	"os"
	"testing"
)

const eps = 1E-4

type readSTLTest struct {
	filename string
	count    int
	t        []STLTriangle
}

var readSTLTests = []readSTLTest{
	{"data/cylinder.bin.stl", 326, nil},
	{
		"data/cylinder.stl",
		326,
		[]STLTriangle{
			{
				[3]float32{0, 0, 0},
				[3]STLPoint{
					{-7.708244e-01, -3.846672e+00, 5.378669e+00},
					{-1.548386e+00, -3.683723e+00, 4.516774e+00},
					{-1.530743e+00, -3.695526e+00, 4.194193e+00},
				},
			},
		},
	},
}

func STLCoordEqual(a, b float32) bool {
	if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
		return false
	}
	return math.Abs(float64(a)-float64(b)) < eps
}

func STLPointEqual(p1, p2 [3]float32) bool {
	for i := 0; i < 3; i++ {
		if !STLCoordEqual(p1[i], p2[i]) {
			return false
		}
	}
	return true
}

func STLEqual(t1, t2 STLTriangle) bool {
	return STLPointEqual(t1.n, t2.n) &&
		STLPointEqual(t1.v[0], t2.v[0]) &&
		STLPointEqual(t1.v[1], t2.v[1]) &&
		STLPointEqual(t1.v[2], t2.v[2])
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
			t.Fatalf("Wrong number of triangles. Expected: %d, got: %d", test.count, len(stl))
		}
		for i, tr := range test.t {
			if !STLEqual(tr, stl[i]) {
				t.Fatalf("Triangle #%d, want: %v, got: %v", i, tr, stl[i])
			}
		}
	}
}
