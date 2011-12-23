package main

type DisjoinSet struct {
	parent []int
	rank   []int
}

func NewDisjoinSet() *DisjoinSet {
	return &DisjoinSet{}
}

func (s *DisjoinSet) Make() int {
	s.parent = append(s.parent, -1)
	s.rank = append(s.rank, 0)
	return len(s.parent) - 1
}

func (s *DisjoinSet) Find(x int) int {
	if s.parent[x] == -1 {
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
