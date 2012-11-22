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

func IsBoundary(vol BoolVoxelVolume, node g3.Node) bool {
	if !vol.Get(node[0], node[1], node[2]) {
		return false
	}
	for _, v := range g3.AdjNodes6 {
		cur := node.Add(v)
		if !vol.Get(cur[0], cur[1], cur[2]) {
			return true
		}
	}
	return false
}
