package main

type Location [2]int64

type LocSet interface {
	Add(loc Location)
	Has(loc Location) bool
	All() []Location
	Clear()
}

type locSet struct {
	a    []int
	l    []Location
	min  [2]int64
	size [2]int64
	b    int
}

func NewLocSet(min [2]int64, size [2]int64) LocSet {
	return &locSet{a: make([]int, size[0]*size[1]), min: min, size: size, b: 1}
}

func (s *locSet) Clear() {
	s.b++
	s.l = s.l[:0]
}

func (s *locSet) index(loc Location) int {
	loc[0] = loc[0] - s.min[0]
	loc[1] = loc[1] - s.min[1]
	return int(loc[0] + loc[1]*s.size[0])
}

func (s *locSet) Add(loc Location) {
	if !s.Has(loc) {
		s.a[s.index(loc)] = s.b
		s.l = append(s.l, loc)
	}
}

func (s *locSet) Has(loc Location) bool {
	return s.a[s.index(loc)] == s.b
}

func (s *locSet) All() (res []Location) {
	res = make([]Location, len(s.l))
	copy(res, s.l)
	return
}
