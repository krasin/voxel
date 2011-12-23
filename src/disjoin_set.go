package main

// DisjoinSet tracks and merges clusters.
// It might be useful for finding graph components and
// related problems. See
// http://en.wikipedia.org/wiki/Union_find for more details.
type DisjoinSet struct {
	parent []int
	rank   []int
}

// NewDisjoinSet creates an empty DisjoinSet.
func NewDisjoinSet() *DisjoinSet {
	return &DisjoinSet{}
}

// Make makes the new cluster.
func (s *DisjoinSet) Make() int {
	x := len(s.parent)
	s.parent = append(s.parent, x)
	s.rank = append(s.rank, 0)
	return x
}

// Find returns the root of cluster #x
func (s *DisjoinSet) Find(x int) int {
	if s.parent[x] != x {
		s.parent[x] = s.Find(s.parent[x])
	}
	return s.parent[x]
}

// Join merges two clusters.
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
