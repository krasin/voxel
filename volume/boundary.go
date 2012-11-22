package volume

import "github.com/krasin/g3"

type BoolVoxelVolume interface {
	Get(x, y, z int) bool
	XLen() int
	YLen() int
	ZLen() int
}

func Normal(vol BoolVoxelVolume, node g3.Node) g3.Vector {
	var p g3.Node

	for _, vec := range g3.AdjNodes26 {
		cur := node.Add(vec)
		if !vol.Get(cur[0], cur[1], cur[2]) {
			continue
		}
		p.Sub(vec)
	}
	if p.IsZero() {
		return g3.Vector{1, 0, 0}
	}
	return p.Vector().Normalize()
}

func IsBoundary(vol BoolVoxelVolume, x, y, z int) bool {
	if !vol.Get(x, y, z) {
		return false
	}
	if x == 0 || x == vol.XLen()-1 || y == 0 || y == vol.YLen()-1 || z == 0 || z == vol.ZLen()-1 {
		return true
	}
	return !(vol.Get(x-1, y, z) && vol.Get(x+1, y, z) &&
		vol.Get(x, y-1, z) && vol.Get(x, y+1, z) &&
		vol.Get(x, y, z-1) && vol.Get(x, y, z+1))
}
