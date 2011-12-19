package main

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
}

func TestSpread3(t *testing.T) {
	for testInd, test := range spread3tests {
		got := spread3(test.num)
		if got != test.want {
			t.Errorf("test #%d: spread3(%d): want %d, got %d", testInd, test.num, test.want, got)
		}
	}
}

type point2htest struct {
	p               Point16
	h               int
	skipReverseTest bool
}

var point2htests = []point2htest{
	{Point16{0, 0, 0}, 0, false},
	{Point16{1, 1, 1}, 0x421, false},
	{Point16{32, 32, 32}, 0, true},
	{Point16{31, 31, 31}, (1 << 15) - 1, false},
	{Point16{15, 31, 4}, 0x3FE4, false},
}

func TestPoint2h(t *testing.T) {
	for testInd, test := range point2htests {
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
