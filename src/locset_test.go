package main

import (
	"testing"
)

const (
	SetAdd   = 0
	SetCheck = 1
	SetClear = 2
)

type locSetTestAction struct {
	Action int
	Loc    Location
	Res    bool
}

type locSetTest struct {
	Size    int64
	Actions []locSetTestAction
}

var locSetTests = []locSetTest{
	{
		Size: 1,
		Actions: []locSetTestAction{
			{SetCheck, Location{0, 0}, false},
			{SetAdd, Location{0, 0}, true},
			{SetCheck, Location{0, 0}, true},
		},
	},
	{
		Size: 2,
		Actions: []locSetTestAction{
			{SetCheck, Location{0, 0}, false},
			{SetCheck, Location{1, 0}, false},
			{SetAdd, Location{0, 0}, true},
			{SetCheck, Location{0, 0}, true},
			{SetCheck, Location{1, 0}, false},
			{Action: SetClear},
			{SetCheck, Location{0, 0}, false},
			{SetCheck, Location{1, 0}, false},
		},
	},
}

func TestLocSet(t *testing.T) {
	for testInd, test := range locSetTests {
		s := NewLocSet([2]int64{}, [2]int64{test.Size, test.Size})
		for actionInd, action := range test.Actions {
			switch action.Action {
			case SetAdd:
				s.Add(action.Loc)
			case SetCheck:
				res := s.Has(action.Loc)
				if res != action.Res {
					t.Errorf("[Test %d, action %d]: check failed. Want: %v, got: %v",
						testInd, actionInd, action.Res, res)
				}
			case SetClear:
				s.Clear()
			}
		}
	}
}
