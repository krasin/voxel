package volume

import (
	"testing"
)

type spread3test struct {
	num  byte
	want int
}

var spread3tests = []spread3test{
	{0, 0},
	{1, 1},
	{3, 9},
	{255, 0x249249},
	{0xAA, 0x208208},
	{0x0F, 0x249},
	{0xDB, 0x241209},
}

func TestSpread3(t *testing.T) {
	for testInd, test := range spread3tests {
		got := spread3(test.num)
		if got != test.want {
			t.Errorf("test #%d: spread3(%d): want %d (0x%x), got %d (0x%x)", testInd, test.num, test.want, test.want, got, got)
		}
		gotNum := join3(test.want)
		if gotNum != test.num {
			t.Errorf("test #%d: join3(%d): want %d, got %d", testInd, test.want, test.num, gotNum)
		}
	}
}

type point2hTest struct {
	p               Point16
	h               int
	skipReverseTest bool
}

var point2hTests = []point2hTest{
	{Point16{0, 0, 0}, 0, false},
	{Point16{1, 1, 1}, 0x421, false},
	{Point16{32, 32, 32}, 0, true},
	{Point16{31, 31, 31}, (1 << 15) - 1, false},
	{Point16{15, 31, 4}, 0x3FE4, false},
}

func TestPoint2h(t *testing.T) {
	for testInd, test := range point2hTests {
		gotH := point2h(test.p)
		if gotH != test.h {
			t.Errorf("test #%d: point2h(%v): want %d, got %d", testInd, test.p, test.h, gotH)
		}
		if test.skipReverseTest {
			continue
		}
		gotP := h2point(test.h)
		if gotP != test.p {
			t.Errorf("test #%d: h2point(%d): want %v, got %v", testInd, test.h, test.p, gotP)
		}
	}
}

type point2kTest struct {
	p Point16
	k int
}

var point2kTests = []point2kTest{
	{Point16{0, 0, 0}, 0},
	{Point16{32, 32, 32}, 7},
	{Point16{0xFF << 5, 0xFF << 5, 0xFF << 5}, 0xFFFFFF},
	{Point16{0, 0, 0xFF << 5}, 0x249249},
	{Point16{103 << 5, 22 << 5, 12 << 5}, 0x1223F4},
}

func TestPoint2k(t *testing.T) {
	for testInd, test := range point2kTests {
		gotK := point2k(test.p)
		if gotK != test.k {
			t.Errorf("test #%d: point2k(%v): want %d (0x%x), got %d (0x%x)", testInd, test.p, test.k, test.k, gotK, gotK)
		}
		gotP := k2point(test.k)
		if gotP != test.p {
			t.Errorf("test #%d: k2point(%d): want %v, got %v", testInd, test.k, test.p, gotP)
		}
	}
}

type point2keyTest struct {
	p   Point16
	key uint64
}

var point2keyTests = []point2keyTest{
	{Point16{0, 0, 0}, 0},
	{Point16{1, 1, 1}, 0x421},
	{Point16{32, 32, 32}, 0x38000},
	{Point16{33, 33, 33}, 0x38421},
	{Point16{511, 0, 1}, 0x4927C01},
}

func TestPoint2Key(t *testing.T) {
	for testInd, test := range point2keyTests {
		got := point2key(test.p)
		if got != test.key {
			t.Errorf("test #%d: point2key(%v): want %d (0x%x), got %d (0x%x)", testInd, test.p, test.key, test.key, got, got)
		}
		gotP := key2point(test.key)
		if gotP != test.p {
			t.Errorf("test #%d: key2point(%d): want %v, got %v", testInd, test.key, test.p, gotP)
		}
	}
}
