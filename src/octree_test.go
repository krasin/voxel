package main

import (
	"testing"
)

const (
	oSet          = 1
	oGet          = 2
	oSetAllFilled = 3
)

type octreeTest struct {
	N      int
	Action []octreeAction
}

type octreeAction struct {
	Op int
	P  [3]int
	V  uint16
}

var octreeTests = []octreeTest{
	{
		512,
		[]octreeAction{
			{oGet, [3]int{0, 0, 0}, 0},
			{oSet, [3]int{0, 0, 0}, 1},
			{oSet, [3]int{0, 0, 1}, 2},
			{oSet, [3]int{0, 1, 1}, 3},
			{oSet, [3]int{1, 1, 1}, 4},
			{oSet, [3]int{511, 0, 0}, 5},
			{oSet, [3]int{511, 511, 511}, 6},
			{oSet, [3]int{511, 511, 1}, 7},
			// Set value to all filled cells.
			{oSetAllFilled, [3]int{}, 10},
			{oGet, [3]int{0, 0, 0}, 10},
			{oGet, [3]int{0, 0, 1}, 10},
			{oGet, [3]int{0, 1, 1}, 10},
			{oGet, [3]int{1, 1, 1}, 10},
			{oGet, [3]int{511, 0, 0}, 10},
			{oGet, [3]int{511, 511, 511}, 10},
			{oGet, [3]int{511, 511, 1}, 10},
		},
	},
}

func TestOctree(t *testing.T) {
	for testInd, test := range octreeTests {
		var tree Uint16Volume = NewOctree(test.N)
		for actInd, act := range test.Action {
			failed := false
			switch act.Op {
			case oSetAllFilled:
				tree.SetAllFilled(act.V)
			case oSet:
				tree.Set(act.P[0], act.P[1], act.P[2], act.V)
				fallthrough
			case oGet:
				v := tree.GetV(act.P[0], act.P[1], act.P[2])
				if v != act.V {
					t.Errorf("test #%d, action #%d: tree%v: want %d, got %d\n", testInd, actInd, act.P, act.V, v)
					failed = true
				}
			default:
				t.Fatalf("test #%d, action #%d: unknown Op: %d\n", testInd, actInd, act.P)
			}
			if failed {
				break
			}
		}
	}
}
