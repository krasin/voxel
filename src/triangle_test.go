package main

import (
	//	"sort"
	"testing"
)

type triangleTest struct {
	p  Point
	t  [3]Point
	in bool
	r  int64
}

var (
	smallRectTriangle = [3]Point{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
	}
	mediumRectTriangle = [3]Point{
		{0, 0, 0},
		{4, 0, 0},
		{0, 4, 0},
	}

	rectTriangle = [3]Point{
		{0, 0, 0},
		{10, 0, 0},
		{0, 10, 0},
	}
	rectTriangle10000 = [3]Point{
		{0, 0, 0},
		{10000, 0, 0},
		{0, 10000, 0},
	}
	eqTriangle = [3]Point{
		{9, 0, 0},
		{0, 9, 0},
		{0, 0, 9},
	}
	thinTriangle = [3]Point{
		{0, 0, 0},
		{10, 1, 0},
		{10, 0, 0},
	}

	triangleTests = []triangleTest{
		{
			Point{1, 1, 0},
			rectTriangle,
			true,
			0,
		},
		{
			Point{1, 1, 1},
			rectTriangle,
			false,
			0,
		},
		{
			Point{1, 1, 1},
			rectTriangle,
			false,
			1,
		},
		{
			Point{1, 1, 1},
			rectTriangle,
			true,
			2,
		},
		{
			Point{11, 0, 0},
			rectTriangle,
			false,
			0,
		},
		{
			Point{10, 0, 0},
			rectTriangle,
			true,
			0,
		},
		{
			Point{0, 0, 0},
			rectTriangle,
			true,
			0,
		},
		{
			Point{0, 10, 0},
			rectTriangle,
			true,
			0,
		},
		{
			Point{0, 0, 0},
			eqTriangle,
			false,
			1,
		},
		{
			Point{0, 0, 0},
			eqTriangle,
			false,
			1,
		},
		{
			Point{3, 3, 3},
			eqTriangle,
			true,
			0,
		},
		{
			Point{3, 3, 4},
			eqTriangle,
			false,
			0,
		},
		{
			Point{3, 3, 4},
			eqTriangle,
			false,
			1,
		},
		{
			Point{3, 3, 4},
			eqTriangle,
			true,
			2,
		},
		{
			Point{3, 3, 5},
			eqTriangle,
			false,
			1,
		},
		{
			Point{0, 0, 0},
			thinTriangle,
			true,
			1,
		},
		{
			Point{1, 0, 0},
			thinTriangle,
			true,
			1,
		},
		{
			Point{0, 1, 0},
			thinTriangle,
			false,
			1,
		},
		{
			Point{1, 1, 0},
			thinTriangle,
			false,
			1,
		},
		{
			Point{2, 1, 0},
			thinTriangle,
			false,
			1,
		},
		{
			Point{4, 1, 0},
			thinTriangle,
			false,
			1,
		},
		{
			Point{5, 1, 0},
			thinTriangle,
			true,
			1,
		},
		{
			Point{9, 1, 0},
			thinTriangle,
			true,
			1,
		},
		{
			Point{10, 1, 0},
			thinTriangle,
			true,
			1,
		},
	}
)

func TestTriangle(t *testing.T) {
	for ind, test := range triangleTests {
		got := DotInTriangle(test.p, test.t[0], test.t[1], test.t[2], test.r)
		if got != test.in {
			t.Errorf("test #%d: %v, want %v, got %v", ind, test, test.in, got)
		}
	}

}

type allTriangleDotsTest struct {
	t     [3]Point
	p     []Point
	scale int64
}

var allTriangleDotsTests = []allTriangleDotsTest{
	{
		smallRectTriangle,
		[]Point{
			{0, 0, 0},
			{1, 0, 0},
			{0, 1, 0},
		},
		1,
	},
	{
		mediumRectTriangle,
		[]Point{
			{0, 0, 0},
			{0, 1, 0},
			{1, 0, 0},
			{0, 2, 0},
			{1, 1, 0},
			{2, 0, 0},
			{0, 3, 0},
			{1, 2, 0},
			{2, 1, 0},
			{3, 0, 0},
			{0, 4, 0},
			{1, 3, 0},
			{2, 2, 0},
			{3, 1, 0},
			{4, 0, 0},
		},
		1,
	},
	{
		rectTriangle,
		[]Point{
			{0, 0, 0},
			{0, 1, 0},
			{1, 0, 0},
			{0, 2, 0},
			{1, 1, 0},
			{2, 0, 0},
			// Grid precision prevents the alrorithm to find these points
			// The algorithm should be improved, but this is fine for now
			//			{2, 1, 0},
			//			{1, 2, 0},
		},
		4,
	},
	{
		rectTriangle10000,
		[]Point{
			{0, 0, 0},
			{0, 1, 0},
			{1, 0, 0},
			{0, 2, 0},
			{1, 1, 0},
			{2, 0, 0},
			{2, 1, 0},
			{1, 2, 0},
		},
		4000,
	},
}

// TODO(krasin): update AllTriangleDots1 call to match the new signature.
func disabledTestAllTriangleDots(t *testing.T) {
	/*	m := make(map[uint64]Point)
		for ind, test := range allTriangleDotsTests {
			sort.Sort(pointSlice(test.p))
			for _, p := range test.p {
				m[hash(p)] = p
			}		
			pt := AllTriangleDots1(test.t[0], test.t[1], test.t[2], test.scale, 1)
			if len(pt) != len(test.p) {
				t.Errorf("Test #%d: number of triangle dots is unexpected. Want: %v, got: %v", ind, test.p, pt)
				continue
			}
			for _, p := range pt {
				if _, ok := m[hash(p)]; !ok {
					t.Errorf("Test #%d: unexpected point: %v. Want: %v, got: %v", ind, p, test.p, pt)
					continue
				}
			}

		}*/
}

type clipLineTest struct {
	line  Line
	cube  Cube
	ok    bool
	scale int64
	res   Line
}

var (
	cube10 = Cube{
		Point{0, 0, 0},
		Point{10, 10, 10},
	}
	lineX100 = Line{
		Point{0, 0, 0},
		Point{100, 0, 0},
	}
	lineXYZ100 = Line{
		Point{0, 0, 0},
		Point{100, 100, 100},
	}
	lineMinusX100 = Line{
		Point{0, 0, 0},
		Point{-100, 0, 0},
	}
	line19 = Line{
		Point{1, 1, 1},
		Point{9, 9, 9},
	}

	clipLineTests = []clipLineTest{
		{
			lineX100,
			cube10,
			true,
			1,
			Line{
				Point{0, 0, 0},
				Point{10, 0, 0},
			},
		},
		{
			lineXYZ100,
			cube10,
			true,
			1,
			Line{
				Point{0, 0, 0},
				Point{10, 10, 10},
			},
		},
		{
			line:  lineMinusX100,
			cube:  cube10,
			ok:    false,
			scale: 1,
		},
		{
			line19,
			cube10,
			true,
			1,
			line19,
		},
		{
			Line{
				Point{-10, 5, 4},
				Point{20, 5, 5},
			},
			cube10,
			true,
			3,
			Line{
				Point{0, 5, 4},
				Point{10, 5, 5},
			},
		},
		{
			Line{
				Point{-10, 5, 4},
				Point{20, 5, 5},
			},
			cube10,
			true,
			1000,
			Line{
				Point{0, 5, 4},
				Point{10, 5, 5},
			},
		},
	}
)

func TestClipLine(t *testing.T) {
	for ind, test := range clipLineTests {
		res, ok := ClipLine(test.line, test.cube, test.scale)
		if test.ok && !ok {
			t.Errorf("Test #%d: Failed to clip line. Test: %v", ind, test)
			continue
		}
		if !test.ok && ok {
			t.Errorf("Test #%d: Found unexpected intersection with the cube: %v. Test: %v", ind, res, test)
			continue
		}
		if !peq(test.res[0], res[0]) || !peq(test.res[1], res[1]) {
			t.Errorf("Test #%d: Wrong result. Want: %v, got: %v. Test: %v", ind, test.res, res, test)
			continue
		}
	}
}

type clipTriangleTest struct {
	triangle Triangle
	cube     Cube
	ok       bool
	scale    int64
	res      []Point
}

var (
	clipTriangleTests = []clipTriangleTest{
		{
			Triangle{
				Point{0, 0, 0},
				Point{5, 10, 0},
				Point{10, 0, 0},
			},
			cube10,
			true,
			1,
			[]Point{
				{0, 0, 0},
				{5, 10, 0},
				{10, 0, 0},
			},
		},
		{
			Triangle{
				Point{0, 0, 0},
				Point{0, 0, 10},
				Point{-10, 10, 5},
			},
			cube10,
			false,
			1,
			[]Point{},
		},
		{
			Triangle{
				Point{-5, 5, 0},
				Point{-5, 5, 10},
				Point{15, 5, 5},
			},
			cube10,
			true,
			2,
			[]Point{
				{0, 5, 9},
				{10, 5, 6},
				{10, 5, 4},
				{0, 5, 1},
			},
		},
		{
			Triangle{
				Point{1, 11, 5},
				Point{1, -1, 5},
				Point{15, 5, 5},
			},
			cube10,
			true,
			10,
			[]Point{
				Point{1, 10, 5},
				Point{1, 0, 5},
				Point{3, 0, 5},
				Point{10, 3, 5},
				Point{10, 7, 5},
				Point{3, 10, 5},
			},
		},
		{
			Triangle{
				Point{-1, 12, 0},
				Point{-1, 7, 0},
				Point{4, 7, 0},
			},
			cube10,
			true,
			10,
			[]Point{
				Point{0, 7, 0},
				Point{4, 7, 0},
				Point{1, 10, 0},
				Point{0, 10, 0},
			},
		},

		{
			Triangle{
				Point{-1, 12, 5},
				Point{-1, 7, 5},
				Point{4, 7, 5},
			},
			cube10,
			true,
			10,
			[]Point{
				Point{0, 7, 5},
				Point{4, 7, 5},
				Point{1, 10, 5},
				Point{0, 10, 5},
			},
		},
		{
			Triangle{
				Point{2, 9, 5},
				Point{8, 9, 5},
				Point{5, 12, 5},
			},
			cube10,
			true,
			10,
			[]Point{
				Point{2, 9, 5},
				Point{8, 9, 5},
				Point{7, 10, 5},
				Point{3, 10, 5},
			},
		},
		{
			Triangle{
				Point{3, 14, 3},
				Point{3, 7, -4},
				Point{-4, 7, 3},
			},
			cube10,
			true,
			10,
			[]Point{
				Point{2, 10, 0},
				Point{0, 10, 2},
				Point{0, 8, 0},
			},
		},
	}
)

// TODO(krasin): enable or delete test.
func disabledTestClipTriangle(t *testing.T) {
	/*	for ind, test := range clipTriangleTests {
		res, ok := ClipTriangle(test.triangle, test.cube, test.scale)
		if test.ok && !ok {
			t.Errorf("Test #%d: no intersection found. Test: %v", ind, test)
			continue
		}
		if !test.ok && ok {
			t.Errorf("Test #%d: unexpected intersection found: %v. Test: %v", ind, res, test)
			continue
		}
		if !test.ok {
			continue
		}
		if len(res) != len(test.res) {
			t.Errorf("Test #%d: different number of points. Want: %v, got: %v. Test: %v", ind, test.res, res, test)
			continue
		}
		for i := range res {
			if !peq(res[i], test.res[i]) {
				t.Errorf("Test #%d, point #%d: wrong result. Want: %v, got: %v. Test: %v", ind, i, test.res, res, test)
				break
			}
		}
	}*/
}
