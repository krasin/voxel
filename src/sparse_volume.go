package main

type Point16 [3]uint16

type SparseVolume struct {
	// Size of leaf cube
	lh uint

	// Depth of octree
	lk uint

	// 2 * lh for speed purposes
	lh2 uint

	// 3 * lh for speed purposes
	lh3 uint

	// (1 << lh) - 1 for speed purposes
	masklh uint64

	// (1 << (3*lh))-1 for speed purposes
	mask3lh uint64
}

func (v *SparseVolume) point2key(p Point16) uint64 {
	return uint64(v.point2k(p))<<v.lh3 + uint64(v.point2h(p))
}

func (v *SparseVolume) point2k(p Point16) int {
	return (spread3(byte(p[0]>>v.lh)) << 2) + (spread3(byte(p[1]>>v.lh)) << 1) + spread3(byte(p[2])>>v.lh)
}

func (v *SparseVolume) k2point(k int) (p Point16) {
	panic("k2point not implemented")
}

func (v *SparseVolume) point2h(p Point16) int {
	return ((int(p[0]) & int(v.masklh)) << v.lh2) + ((int(p[1]) & int(v.masklh)) << v.lh) + (int(p[2]) & int(v.masklh))
}

func (v *SparseVolume) h2point(h int) (p Point16) {
	p[0] = uint16(h >> v.lh2)
	p[1] = uint16((h >> v.lh) & int(v.masklh))
	p[2] = uint16(h & int(v.masklh))
	return
}

func spread3(b byte) (x int) {
	x = int(b)
	x = ((x & 0xF0) << 8) | (x & 0x0F)
	x = ((x & 0xC00C) << 4) | (x & 0x3003)
	x = ((x & 0x82082) << 2) | (x & 0x41041)
	return
}

func (v *SparseVolume) key2h(key uint64) int {
	return int(key & v.mask3lh)
}

func (v *SparseVolume) key2k(key uint64) int {
	return int(key >> v.lh3)
}

func (v *SparseVolume) key2point(key uint64) (p Point16) {
	ph := v.h2point(v.key2h(key))
	pk := v.k2point(v.key2k(key))
	p[0] = (pk[0] << v.lh) | ph[0]
	p[1] = (pk[1] << v.lh) | ph[1]
	p[2] = (pk[2] << v.lh) | ph[2]
	return
}
