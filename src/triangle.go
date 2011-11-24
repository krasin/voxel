package main

import (
	"big"
)

type Point [3]int64
type Vector [3]int64
type Triangle [3]Point

func NewVector(p1, p2 Point) Vector {
	return Vector{p2[0] - p1[0], p2[1] - p1[1], p2[2] - p1[2]}
}

func VectorProduct(a, b Vector) Vector {
	return Vector{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func ScalarProduct(a, b Vector) (s *big.Int) {
	s = big.NewInt(0)
	for i := 0; i < 3; i++ {
		tmp := big.NewInt(0)
		s.Add(s, tmp.Mul(big.NewInt(a[i]), big.NewInt(b[i])))
	}
	return
}

func DotInPlane(p, a, b, c Point, r int64) bool {
	va := NewVector(c, a)
	vb := NewVector(c, b)
	vc := NewVector(c, p)
	v := VectorProduct(va, vb)

	s := ScalarProduct(vc, v)
	s.Mul(s, s)
	s.Mul(s, big.NewInt(4))

	r2 := big.NewInt(r)
	r2.Mul(r2, r2)

	v2 := ScalarProduct(v, v)

	t := big.NewInt(0)
	t.Mul(r2, v2)

	return t.Cmp(s) >= 0
}

func len2(v Vector) (s *big.Int) {
	s = big.NewInt(0)
	for i := 0; i < 3; i++ {
		tmp := big.NewInt(v[i])
		tmp.Mul(tmp, tmp)
		s.Add(s, tmp)
	}
	return
}

func sameSide(p, a, b, c Point, r int64) bool {
	ab := NewVector(a, b)
	ac := NewVector(a, c)
	ap := NewVector(a, p)
	v1 := VectorProduct(ab, ac)
	v2 := VectorProduct(ab, ap)
	s := ScalarProduct(v1, v2)
	if s.Cmp(big.NewInt(0)) >= 0 {
		return true
	}
	h2 := len2(v2)
	h2.Mul(h2, big.NewInt(4))
	r2 := big.NewInt(r)
	r2.Mul(r2, r2)
	r2.Mul(r2, len2(ab))
	return r2.Cmp(h2) >= 0
}

func inplaneDotInTriangle(p, a, b, c Point, r int64) bool {
	return sameSide(p, a, b, c, r) &&
		sameSide(p, b, c, a, r) &&
		sameSide(p, c, a, b, r)
}

func DotInTriangle(p, a, b, c Point, r int64) bool {
	return DotInPlane(p, a, b, c, r) && inplaneDotInTriangle(p, a, b, c, r)
}

func hash(p Point) uint64 {
	return (uint64(p[0]) << 42) + (uint64(p[1] << 21)) + uint64(p[2])
}

func adjacent(p Point) (res []Point) {
	res = make([]Point, 8)[0:0]
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			for dz := -1; dz <= 1; dz++ {
				if dx == 0 && dy == 0 && dz == 0 {
					continue
				}
				res = append(res, Point{
					p[0] + int64(dx),
					p[1] + int64(dy),
					p[2] + int64(dz),
				})
			}
		}
	}
	return
}

func toGrid(p Point, scale int64) Point {
	return Point{p[0] / scale, p[1] / scale, p[2] / scale}
}

func fromGrid(p Point, scale int64) Point {
	return Point{p[0]*scale + scale/2, p[1]*scale + scale/2, p[2]*scale + scale/2}
}

func AllTriangleDots(a, b, c Point, scale, r int64) (res []Point) {
	q := []Point{toGrid(a, scale), toGrid(b, scale), toGrid(c, scale)}
	var q2 []Point
	m := make(map[uint64]Point)
	m[hash(a)] = a
	m[hash(b)] = b
	m[hash(c)] = c
	for len(q) > 0 {
		q, q2 = q2[0:0], q
		for _, p := range q2 {
			for _, p2 := range adjacent(p) {
				if _, ok := m[hash(p2)]; ok {
					continue
				}

				if !DotInTriangle(fromGrid(p2, scale), a, b, c, r) {
					continue
				}
				m[hash(p2)] = p2
				q = append(q, p2)
			}
		}
	}
	res = make([]Point, len(m))[0:0]
	for _, p := range m {
		res = append(res, p)
	}
	return
}
