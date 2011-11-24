package main

import (
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
	t [3]Point
	p []Point
}

var allTriangleDotsTests = []allTriangleDotsTest{
	{
		smallRectTriangle,
		[]Point{
			{0, 0, 0},
			{1, 0, 0},
			{0, 1, 0},
		},
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
	},
}

func TestAllTriangleDots(t *testing.T) {
	m := make(map[uint64]Point)
	for ind, test := range allTriangleDotsTests {
		for _, p := range test.p {
			m[hash(p)] = p
		}
		pt := AllTriangleDots(test.t[0], test.t[1], test.t[2], 1)
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

	}
}
