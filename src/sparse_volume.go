package main

type Point16 [3]uint16

const (
	// Size of leaf cube
	lh = 5

	masklh  = (1 << lh) - 1
	mask3lh = (1 << 3 * lh) - 1
)

type SparseVolume struct {
	// Depth of octree
	lk uint
}

func (v *SparseVolume) point2key(p Point16) uint64 {
	return uint64(point2k(p))<<(3*lh) + uint64(point2h(p))
}

func point2k(p Point16) int {
	return (spread3(byte(p[0]>>lh)) << 2) + (spread3(byte(p[1]>>lh)) << 1) + spread3(byte(p[2]>>lh))
}

func k2point(k int) (p Point16) {
	p[0] = uint16(join3((k>>2)&0x249249)) << lh
	p[1] = uint16(join3((k>>1)&0x249249)) << lh
	p[2] = uint16(join3(k&0x249249)) << lh
	return
}

func point2h(p Point16) int {
	return ((int(p[0]) & masklh) << (2 * lh)) + ((int(p[1]) & masklh) << lh) + (int(p[2]) & masklh)
}

func h2point(h int) (p Point16) {
	p[0] = uint16(h >> (2 * lh))
	p[1] = uint16((h >> lh) & masklh)
	p[2] = uint16(h & masklh)
	return
}

func spread3(b byte) (x int) {
	x = int(b)
	x = ((x & 0xF0) << 8) | (x & 0x0F)
	x = ((x & 0xC00C) << 4) | (x & 0x3003)
	x = ((x & 0x82082) << 2) | (x & 0x41041)
	return
}

func join3(x int) (b byte) {
	x = ((x & 0x208208) >> 2) | (x & 0xDF7DF7)
	x = ((x & 0xC00C0) >> 4) | (x & 0x3FF3F)
	x = ((x & 0xF000) >> 8) | (x & 0x0FFF)
	return byte(x)
}

func key2h(key uint64) int {
	return int(key & mask3lh)
}

func key2k(key uint64) int {
	return int(key >> 3 * lh)
}

func key2point(key uint64) (p Point16) {
	ph := h2point(key2h(key))
	pk := k2point(key2k(key))
	p[0] = (pk[0] << lh) | ph[0]
	p[1] = (pk[1] << lh) | ph[1]
	p[2] = (pk[2] << lh) | ph[2]
	return
}
