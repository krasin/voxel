package main

import (
	"big"
)

type Point [3]int64
type Vector [3]int64

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

	r2 := big.NewInt(r)
	r2.Mul(r2, r2)

	v2 := ScalarProduct(v, v)

	t := big.NewInt(0)
	t.Mul(r2, v2)

	return t.Cmp(s) >= 0
}

func sameSide(p, a, b, c Point) bool {
	ab := NewVector(a, b)
	ac := NewVector(a, c)
	ap := NewVector(a, p)
	v1 := VectorProduct(ab, ac)
	v2 := VectorProduct(ab, ap)
	s := ScalarProduct(v1, v2)
	return s.Cmp(big.NewInt(0)) >= 0
}

func inplaneDotInTriangle(p, a, b, c Point) bool {
	return sameSide(p, a, b, c) &&
		sameSide(p, b, c, a) &&
		sameSide(p, c, a, b)
}

func DotInTriangle(p, a, b, c Point, r int64) bool {
	return DotInPlane(p, a, b, c, r) && inplaneDotInTriangle(p, a, b, c)
}
