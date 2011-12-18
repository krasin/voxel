package main

import (
	"fmt"
	"math"
)

const primeIndexSize = 1 << 28

var octShift = calcOctShift(10)

func calcOctShift(n int) (res []int64) {
	res = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		res[i] = (res[i-1] << 3) + 1
	}
	return
}

type Octree struct {
	N        int
	leafSize int
	mask     int64
	//	p        [][]uint16
	v []uint16
}

func log2(n int64) (res uint) {
	for n > 1 {
		n >>= 1
		res++
	}
	return
}

// N must be power of 2.
func NewOctree(N int) *Octree {
	volume := int64(1) << (3 * log2(int64(N)))
	indexSize := int64(primeIndexSize)
	//	indexSize := int64(1) << log2(volume)
	//	if indexSize > primeIndexSize {
	//		panic("extended part of octree is not implemented")
	//		//		indexSize = primeIndexSize
	//	}
	return &Octree{
		N:        N,
		leafSize: int(volume / indexSize),
		mask:     int64(indexSize - 1),
		//		p:        make([][]uint16, indexSize),
		v: make([]uint16, indexSize),
	}
}

func (t *Octree) internalGet(depth uint, p, base [3]int, l int, index int64) uint16 {
	if index+octShift[depth] > t.mask {
		panic(fmt.Sprintf("Octree.internalGet: extended part of octree is not implemented. depth: %d, t.mask: %d, index: %d, index+octShift[depth]: %d, p: %v", depth, t.mask, index, index+octShift[depth], p))
	}
	//	fmt.Printf("internalGet, depth: %d, p: %v, base: %v, l: %d, index: %d, t.mask: %d, index+octShift[depth]: %d\n", depth, p, base, l, index, t.mask, index+octShift[depth])
	arindex := int((index + octShift[depth]) & t.mask)
	//	if t.p[arindex] != nil {
	//		panic("extended part of octree is not implemented")
	//	}
	if t.v[arindex] == 0 {
		return 0
	}
	if l == 1 || t.v[arindex] != 1 {
		// This is the leaf
		return t.v[arindex] - 2
	}
	// t.v[arindex] == 1 which is the special value means "see subcube"
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
	if x < 0 || y < 0 || z < 0 || x >= t.N || y >= t.N || z >= t.N {
		return 0
	}
	return t.internalGet(0, [3]int{x, y, z}, [3]int{0, 0, 0}, t.N, 0)
}

func (t *Octree) Get(x, y, z int) bool {
	return t.GetV(x, y, z) != 0
}

func (t *Octree) XLen() int {
	return t.N
}

func (t *Octree) YLen() int {
	return t.N
}

func (t *Octree) ZLen() int {
	return t.N
}

func (t *Octree) internalSet(depth uint, p, base [3]int, l int, index int64, v uint16) {
	//	fmt.Printf("internalSet(depth=%d, p=%v, base=%v, l=%d, index=%d, v=%d)", depth, p, base, l, index, v)
	if v >= math.MaxUint16-1 {
		panic("v >= math.MaxUint16-1. These values are reserved")
	}
	if index+octShift[depth] > t.mask {
		panic("extended part of octree is not implemented")
	}
	arindex := int((index + octShift[depth]) & t.mask)
	//	fmt.Printf(", arindex: %d, t.v[arindex] = %d\n", arindex, t.v[arindex])
	//	if t.p[arindex] != nil {
	//		panic("extended part of octree is not implemented")
	//	}
	if l == 1 {
		t.v[arindex] = v + 2
		return
	}
	l >>= 1
	cur := t.v[arindex]
	switch {
	case cur == 0:
		t.v[arindex] = 1
	case cur == 1:
	case cur == v+2:
		// nothing to do
		return
	default:
		// We need to split the cube into 8 smaller cubes and
		// recurse into one of them
		for i := 0; i < 8; i++ {
			nindex := (index << 3) + int64(i)
			if nindex+octShift[depth+1] > t.mask {
				//				fmt.Printf("depth: %d, nindex: %d, depth: %d, octShift[depth+1]: %d, t.mask: %d\n", depth, nindex, depth, octShift[depth+1], t.mask)
				panic("extended part of octree is not implemented")
			}
			arnindex := int((nindex + octShift[depth+1]) & t.mask)
			//			fmt.Printf("t.v[%d] = %d\n", arnindex, cur)
			t.v[arnindex] = cur
		}
	}

	for i := 0; i < 3; i++ {
		index <<= 1
		if p[i] >= base[i]+l {
			index++
			base[i] += l
		}
	}
	t.internalSet(depth+1, p, base, l, index, v)
}

func (t *Octree) Set(x, y, z int, v uint16) {
	t.internalSet(0, [3]int{x, y, z}, [3]int{0, 0, 0}, t.N, 0, v)
}

func (t *Octree) internalSetAllFilled(depth uint, base [3]int, l int, index int64, v uint16) {
	if v >= math.MaxUint16-1 {
		panic("v >= math.MaxUint16-1. These values are reserved")
	}
	if index+octShift[depth] > t.mask {
		panic("extended part of octree is not implemented")
	}
	arindex := int((index + octShift[depth]) & t.mask)
	l >>= 1
	switch t.v[arindex] {
	// This node is empty, nothing to do
	case 0:

	// We need to visit subnodes
	case 1:
		for i := 0; i < 8; i++ {
			t.internalSetAllFilled(
				depth+1,
				[3]int{base[0] + l*(i&4), base[1] + l*(i&2), base[2] + l*(i&1)},
				l,
				(index<<3)+int64(i),
				v)

		}

	// This node is empty (unfortunately, we have 2 states for empty nodes)
	case 2:

	// Child node
	default:
		t.v[arindex] = v + 2
	}
}

func (t *Octree) SetAllFilled(v uint16) {
	t.internalSetAllFilled(0, [3]int{0, 0, 0}, t.N, 0, v)
}

func (t *Octree) internalMapBoundary(depth uint, base [3]int, l int, index int64, f func(x, y, z int)) {
	if index+octShift[depth] > t.mask {
		panic("extended part of octree is not implemented")
	}
	arindex := int((index + octShift[depth]) & t.mask)
	l >>= 1
	switch t.v[arindex] {
	// This node is empty, nothing to do
	case 0:

	// We need to visit subnodes
	case 1:
		for i := 0; i < 8; i++ {
			t.internalMapBoundary(
				depth+1,
				[3]int{base[0] + l*(i&4), base[1] + l*(i&2), base[2] + l*(i&1)},
				l,
				(index<<3)+int64(i),
				f)

		}
	// This node is empty (unfortunately, we have 2 states for empty nodes)
	case 2:

	// Child node
	default:
		if l == 0 {
			if IsBoundary(t, base[0], base[1], base[2]) {
				f(base[0], base[1], base[2])
			}
			return
		}

		panic(fmt.Sprintf("Visiting edges of the cube is not implemented, l=%d, t.v[arindex]: %d", l, t.v[arindex]))
	}
}

func (t *Octree) MapBoundary(f func(x, y, z int)) {
	t.internalMapBoundary(0, [3]int{0, 0, 0}, t.N, 0, f)
}
