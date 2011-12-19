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
		has := spread3(test.num)
		if has != test.want {
			t.Errorf("test #%d: spread3(%d) returned %d, want %d", testInd, test.num, has, test.want)
		}
	}
}
