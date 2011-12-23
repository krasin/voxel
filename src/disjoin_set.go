package main

type DisjoinSet struct {
	parent []int
	rank   []int
}

func NewDisjoinSet() *DisjoinSet {
	return &DisjoinSet{}
}

func (s *DisjoinSet) Make() int {
	x := len(s.parent)
	s.parent = append(s.parent, x)
	s.rank = append(s.rank, 0)
	return x
}

func (s *DisjoinSet) Find(x int) int {
	if s.parent[x] == x {
		return x
	}
	return s.Find(s.parent[x])
}

func (s *DisjoinSet) Join(x, y int) {
	xRoot := s.Find(x)
	yRoot := s.Find(y)

	switch {
	case s.rank[xRoot] < s.rank[yRoot]:
		s.parent[xRoot] = yRoot
	case s.rank[xRoot] > s.rank[yRoot]:
		s.parent[yRoot] = xRoot
	default:
		s.rank[xRoot]++
		s.parent[yRoot] = xRoot
	}
}
