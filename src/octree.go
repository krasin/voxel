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

func (t *Octree) internalGet(depth uint, p [3]int, base [3]int, l int, index int64) uint16 {
	if t.p[int(index&t.mask)] != nil {
		panic("extended part of octree is not implemented")
	}
	if t.v[int(index&t.mask)] == 0 {
		return 0
	}
	if l == 1 {
		return t.v[int(index&t.mask)] - 1
	}
	// We need to check the neighbour if it has the value of 0 or not.
	// If no, we need to go deeper.
	// This is because we reuse the cell after an expansion.
	another := index + (1 << (depth * 3))
	if another > t.mask || t.v[int(another&t.mask)] == 0 {
		// This is the leaf
		return t.v[int(index&t.mask)] - 1
	}
	l >>= 1
	for i := 0; i < 3; i++ {
		index <<= 1
		if p[i] >= base[i]+l {
			index++
			base[i] += l
		}
	}

	return t.internalGet(depth+1, p, base, l, index)
}

func (t *Octree) GetV(x, y, z int) uint16 {
	return t.internalGet(1, [3]int{x, y, z}, [3]int{0, 0, 0}, t.N, 0)
}
