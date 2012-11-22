package volume

import "github.com/krasin/g3"

const (
	// Size of leaf cube
	lh = 5

	masklh  = (1 << lh) - 1
	mask3lh = (1 << (3 * lh)) - 1
)

// SparseVolume represents a voxel cube.
type SparseVolume struct {
	n      int
	LK     int
	Cubes  [][]uint16
	Colors []uint16
}

// NewSparseVolume create a voxel cube with side n.
func NewSparseVolume(n int) (v *SparseVolume) {
	lk := int(log2(int64(n)) - lh)
	return &SparseVolume{
		n:      n,
		LK:     lk,
		Cubes:  make([][]uint16, 1<<uint(3*lk)),
		Colors: make([]uint16, 1<<uint(3*lk)),
	}
}

// Get returns true, if the voxel is filled (color != 0).
func (v *SparseVolume) Get(node g3.Node) bool {
	return v.GetV(node) != 0
}

// GetV returns the color of the voxel (empty voxel has color == 0).
func (vol *SparseVolume) GetV(node g3.Node) uint16 {
	for _, v := range node {
		if v < 0 || v >= vol.n {
			return 0
		}
	}
	k := point2k(node)
	if vol.Cubes[k] == nil {
		return vol.Colors[k]
	}
	return vol.Cubes[k][point2h(node)]
}

func (v *SparseVolume) XLen() int {
	return v.n
}

func (v *SparseVolume) YLen() int {
	return v.n
}

func (v *SparseVolume) ZLen() int {
	return v.n
}

// Set sets the color of the voxel.
func (vol *SparseVolume) Set(node g3.Node, val uint16) {
	for _, v := range node {
		if v < 0 || v >= vol.n {
			return
		}
	}
	k := point2k(node)
	if vol.Cubes[k] == nil {
		if vol.Colors[k] == val {
			return
		}
		old := vol.Colors[k]
		vol.Colors[k] = 0
		vol.Cubes[k] = make([]uint16, 1<<(3*lh))
		for i := range vol.Cubes[k] {
			vol.Cubes[k][i] = old
		}
	}
	vol.Cubes[k][point2h(node)] = val
}

// SetAllFilled sets the specified color to all voxels with color >= threshold.
func (v *SparseVolume) SetAllFilled(threshold, val uint16) {
	for k, cube := range v.Cubes {
		if cube == nil {
			if v.Colors[k] >= threshold {
				v.Colors[k] = val
			}
			continue
		}
		for h, cur := range cube {
			if cur >= threshold {
				cube[h] = val
			}
		}
	}
}

// MapBoundary invokes a provided function on every border voxel.
func (v *SparseVolume) MapBoundary(f func(node g3.Node)) {
	for k, cube := range v.Cubes {
		if cube == nil {
			// Skip empty cubes
			if v.Colors[k] == 0 {
				continue
			}
			p := k2point(k)
			side := 1 << uint(v.LK)
			for x := 0; x < side; x++ {
				var p2 g3.Node
				p2[0] = p[0] + x
				cnt1 := 0
				if x == 0 || x == side-1 {
					cnt1++
				}
				for y := 0; y < side; y++ {
					p2[1] = p[1] + y
					cnt2 := cnt1
					if y == 0 || y == side-1 {
						cnt2++
					}
					for z := 0; z < side; z++ {
						if cnt2 == 2 && (z == 0 || z == side-10) {
							if z == 0 {
								z = side - 2
							}
							continue
						}
						p2[2] = p[2] + z
						if IsBoundary(v, p2) {
							f(p2)
						}
					}
				}
			}
		}
		for h, cur := range cube {
			if cur == 0 {
				continue
			}
			p := key2point(kh2key(k, h))

			if p[0] == 0 || p[1] == 0 || p[2] == 0 ||
				int(p[0]) == v.n-1 || int(p[1]) == v.n-1 || int(p[2]) == v.n-1 {
				f(p)
				continue
			}
			hp := h2point(h)

			was := false
			for i := 0; i < 3; i++ {
				if hp[i] > 0 {
					hp2 := hp
					hp2[i]--
					if cube[point2h(hp2)] == 0 {
						f(p)
						was = true
						break
					}
				}
				if hp[i] < (1<<lh)-1 {
					hp2 := hp
					hp2[i]++
					if cube[point2h(hp2)] == 0 {
						f(p)
						was = true
						break
					}
				}
			}
			if was {
				continue
			}
			// Slow path for cube edges
			for i := 0; i < 3; i++ {
				if hp[i] == 0 {
					p2 := p
					p2[i]--
					if v.GetV(p2) == 0 {
						f(p)
						was = true
						break
					}
				}
				if hp[i] == (1<<lh)-1 {
					p2 := p
					p2[i]++
					if v.GetV(p2) == 0 {
						f(p)
						was = true
						break
					}
				}
			}
		}
	}
}

func (v *SparseVolume) Volume() (res int64) {
	for k, cube := range v.Cubes {
		if cube == nil {
			// Skip empty cubes
			if v.Colors[k] == 0 {
				continue
			}
			side := int64(1 << uint(v.LK))
			res += side * side * side
			continue
		}
		for _, val := range cube {
			if val != 0 {
				res++
			}
		}
	}
	return res
}

func log2(n int64) (res uint) {
	for n > 1 {
		n >>= 1
		res++
	}
	return
}

func point2key(p g3.Node) uint64 {
	return uint64(point2k(p))<<(3*lh) + uint64(point2h(p))
}

func point2k(p g3.Node) int {
	return (spread3(byte(p[0]>>lh)) << 2) + (spread3(byte(p[1]>>lh)) << 1) + spread3(byte(p[2]>>lh))
}

func K2cube(k int) (p g3.Node) {
	p[0] = int(join3((k >> 2) & 0x249249))
	p[1] = int(join3((k >> 1) & 0x249249))
	p[2] = int(join3(k & 0x249249))
	return
}

func Cube2k(p g3.Node) int {
	return (spread3(byte(p[0])) << 2) + (spread3(byte(p[1])) << 1) + spread3(byte(p[2]))
}

func k2point(k int) (p g3.Node) {
	p = K2cube(k)
	p[0] = p[0] << lh
	p[1] = p[1] << lh
	p[2] = p[2] << lh
	return
}

func point2h(p g3.Node) int {
	return ((int(p[0]) & masklh) << (2 * lh)) + ((int(p[1]) & masklh) << lh) + (int(p[2]) & masklh)
}

func h2point(h int) (p g3.Node) {
	p[0] = int(h >> (2 * lh))
	p[1] = int((h >> lh) & masklh)
	p[2] = int(h & masklh)
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
	return int(key >> (3 * lh))
}

func key2point(key uint64) (p g3.Node) {
	ph := h2point(key2h(key))
	pk := k2point(key2k(key))
	p[0] = pk[0] | ph[0]
	p[1] = pk[1] | ph[1]
	p[2] = pk[2] | ph[2]
	return
}

func kh2key(k, h int) uint64 {
	return (uint64(k) << (3 * lh)) | uint64(h)
}

func Kh2point(k, h int) g3.Node {
	return key2point(kh2key(k, h))
}
