package volume

import "github.com/krasin/g3"

type Space interface {
	Get(node g3.Node) bool
	N() int
}

type Uint16Space interface {
	Space
	Set(node g3.Node, v uint16)
	GetV(node g3.Node) uint16
	SetAllFilled(threshold, val uint16)
	MapBoundary(f func(node g3.Node))
	Volume() int64
}

func Normal(vol Space, node g3.Node) g3.Vector {
	var p g3.Node

	for _, vec := range g3.AdjNodes26 {
		cur := node.Add(vec)
		if !vol.Get(cur) {
			continue
		}
		p.Sub(vec)
	}
	if p.IsZero() {
		return g3.Vector{1, 0, 0}
	}
	return p.Vector().Normalize()
}

func IsBoundary(vol Space, node g3.Node) bool {
	if !vol.Get(node) {
		return false
	}
	for _, v := range g3.AdjNodes6 {
		cur := node.Add(v)
		if !vol.Get(cur) {
			return true
		}
	}
	return false
}
