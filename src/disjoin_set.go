package main

type DisjoinSet struct {
	parents []int
}

func NewDisjoinSet() *DisjoinSet {
	return &DisjoinSet{}
}

func (s *DisjoinSet) Make() int {
	s.parents = append(s.parents, -1)
	return len(s.parents) - 1
}

func (s *DisjoinSet) Find(x int) int {
	if s.parents[x] == -1 {
		return x
	}
	return s.Find(s.parents[x])
}

func (s *DisjoinSet) Join(x, y int) {
	xRoot := s.Find(x)
	yRoot := s.Find(y)
	s.parents[yRoot] = xRoot
}
