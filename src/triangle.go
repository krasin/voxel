package main

import (
	"big"
)

type Point [3]int64
type Vector [3]int64
type Triangle [3]Point
type Line [2]Point
type Cube [2]Point

const MaxJ = 10

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

func to(a, n int64) int64 {
	if a >= 0 {
		return (a + (n-1)/2) / n
	}
	return -((-a + n/2) / n)
}

func toGrid(p Point, scale int64) Point {
	return Point{to(p[0], scale), to(p[1], scale), to(p[2], scale)}
}

func scalePoint(p Point, scale int64) Point {
	return Point{p[0] * scale, p[1] * scale, p[2] * scale}
}

func findJ(p1, p2 Point, scale int64) (j uint) {
	for j = 0; j < 31; j++ {
		var r2 int64
		for z := 0; z < 3; z++ {
			diff := int64(p1[z] - p2[z])
			r2 += diff * diff
		}
		if r2 < (int64(scale)*int64(scale))<<(2*j) {
			return j + 2
		}
	}
	panic("unreachable")
}

func peq(p1, p2 Point) bool {
	for z := 0; z < 3; z++ {
		if p1[z] != p2[z] {
			return false
		}
	}
	return true
}

type pointSlice []Point

func (ps pointSlice) Len() int {
	return len(ps)
}

func (ps pointSlice) Less(i, j int) (res bool) {
	for z := 0; z < 3; z++ {
		if ps[i][z] < ps[j][z] {
			return true
		}
		if ps[i][z] > ps[j][z] {
			return false
		}
	}
	return false
}

func (ps pointSlice) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func AddDot(a, b, c Point, scale int64, vol VolumeSetter, i0, i1 int64, j0, j1 uint, last1 Point, color uint16) Point {
	m := j0
	if m < j1 {
		m = j1
	}
	i2 := 1<<m - i0*(1<<(m-j0)) - i1*(1<<(m-j1))
	var p Point
	for z := 0; z < 3; z++ {
		p[z] = int64(i0)*(int64(1)<<uint(m-j0))*a[z] +
			int64(i1)*(int64(1)<<uint(m-j1))*b[z] +
			int64(i2)*c[z]
		p[z] >>= m
	}

	p = toGrid(p, scale)
	vol.Set(int(p[0]), int(p[1]), int(p[2]), color)

	return p
}

func AllTriangleDots(a, b, c Point, scale int64, vol VolumeSetter, color uint16) {
	j0 := findJ(a, c, scale)
	j1 := findJ(a, b, scale)

	m := j0
	if m < j1 {
		m = j1
	}

	for i0 := 0; i0 <= 1<<j0; i0++ {
		var last1 Point
		for i1 := 0; i0*(1<<(m-j0))+i1*(1<<(m-j1)) <= 1<<m; i1++ {
			last1 = AddDot(a, b, c, scale, vol, int64(i0), int64(i1), j0, j1, last1, color)

		}
	}
}

func det3(v0, v1, v2 Vector) int64 {
	return v0[0]*v1[1]*v2[2] + v0[1]*v1[2]*v2[0] + v0[2]*v1[0]*v2[1] -
		v0[0]*v1[2]*v2[1] - v0[1]*v1[0]*v2[2] - v0[2]*v1[1]*v2[0]
}

func MeshVolume(triangles []Triangle, scale int64) (res int64) {
	for _, t := range triangles {
		res += det3(Vector(t[0]), Vector(t[1]), Vector(t[2]))
	}
	return res / 6
}
