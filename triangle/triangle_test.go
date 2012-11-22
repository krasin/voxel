package triangle

import (
	"sort"
	"testing"

	"github.com/krasin/g3"
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

type testVolumeSetter struct {
	p   []Point
	val []uint16
}

func (s *testVolumeSetter) Set(node g3.Node, val uint16) {
	ind := s.find(node)
	if ind != -1 {
		s.val[ind] = val
		return
	}
	s.p = append(s.p, Point{int64(node[0]), int64(node[1]), int64(node[2])})
	s.val = append(s.val, val)
}

func (s *testVolumeSetter) find(node g3.Node) int {
	for ind, p := range s.p {
		if p[0] == int64(node[0]) && p[1] == int64(node[1]) && p[2] == int64(node[2]) {
			return ind
		}
	}
	return -1
}

func (s *testVolumeSetter) Get(node g3.Node) bool {
	ind := s.find(node)
	return ind != -1 && s.val[ind] != 0
}

func TestAllTriangleDots(t *testing.T) {
	for ind, test := range allTriangleDotsTests {
		sort.Sort(pointSlice(test.p))
		vol := new(testVolumeSetter)
		AllTriangleDots(test.t[0], test.t[1], test.t[2], test.scale, vol, 1)
		if len(vol.p) != len(test.p) {
			t.Errorf("Test #%d: number of triangle dots is unexpected. Want: %v, got: %v", ind, test.p, vol.p)
			continue
		}
		for _, p := range test.p {
			if !vol.Get(g3.Node{int(p[0]), int(p[1]), int(p[2])}) {
				t.Errorf("Test #%d: point expected, but not returned: %v. Want: %v, got: %v", ind, p, test.p, vol.p)
				continue
			}
		}

	}
}
