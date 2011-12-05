package main

const primeIndexSize = 1 << 24

type Octree struct {
	N        int
	leafSize int
	mask     int64
	p        [][]uint16
	v        []uint16
}

// N must be power of 2.
func NewOctree(N int) *Octree {
	indexSize := N
	if N > primeIndexSize {
		indexSize = primeIndexSize
	}
	return &Octree{
		N:        N,
		leafSize: N / indexSize,
		mask:     int64(indexSize - 1),
		p:        make([][]uint16, indexSize),
		v:        make([]uint16, indexSize),
	}
}

func step(p, base [3]int, l int, index int64) (nbase [3]int, nl int, nindex int64) {
	nbase, nl, nindex = base, l>>1, index
	for i := 0; i < 3; i++ {
		nindex <<= 1
		if p[i] >= nbase[i]+nl {
			nindex++
			nbase[i] += nl
		}
	}
	return
}

func (t *Octree) GetV(x, y, z int) uint16 {
	p := [3]int{x, y, z}
	base := [3]int{0, 0, 0}
	var index int64
	l := t.N
	count := uint(1)
	base, l, index = step(p, base, l, index)
	if t.p != nil {
		panic("extended part of octree is not implemented")
	}
	cur := t.v[int(index&int64(t.mask))]
	if cur == 0 {
		return 0
	}

	// We need to check the neighbour if it has the value of 0 or not.
	// If no, we need to go deeper.
	// This is because we reuse the cell after an expansion.
	another := index + (1 << (count * 3))
	if another > t.mask || t.v[int(another&t.mask)] == 0 {
		// This is the leaf
		return cur
	}

	// The node is not leaf
	panic("not implemented")
}
