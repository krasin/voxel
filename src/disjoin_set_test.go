package main

import (
	"testing"
)

const (
	dsMake = 0
	dsFind = 1
	dsJoin = 2
)

type disjoinSetAct struct {
	Op int
	X  int
	Y  int
}

type disjoinSetTest struct {
	Act []disjoinSetAct
}

var disjoinSetTests = []disjoinSetTest{
	{
		[]disjoinSetAct{
			{Op: dsMake, X: 0},
			{Op: dsMake, X: 1},
			{Op: dsMake, X: 2},
			{Op: dsFind, X: 0, Y: 0},
			{Op: dsFind, X: 1, Y: 1},
			{Op: dsJoin, X: 0, Y: 1},
			{Op: dsFind, X: 0, Y: 0},
			{Op: dsFind, X: 1, Y: 0},
		},
	},
}

func TestDisjoinSet(t *testing.T) {
	for testInd, test := range disjoinSetTests {
		s := NewDisjoinSet()
		for actInd, act := range test.Act {
			switch act.Op {
			case dsMake:
				x := s.Make()
				if x != act.X {
					t.Errorf("Test #%d, act #%d, Op: dsMake, unexpected result. Want %d, got %d",
						testInd, actInd, x, act.X)
				}
			case dsFind:
				y := s.Find(act.X)
				if y != act.Y {
					t.Errorf("Test #%d, act #%d, Op: dsFind(%d), unexpected result. Want %d, got %d",
						testInd, actInd, act.X, act.Y, y)
				}
			case dsJoin:
				s.Join(act.X, act.Y)
			default:
				t.Fatalf("Test #%d, act #%d, unknown Op: %d", testInd, actInd, act.Op)
			}
		}
	}
}
