package main

import (
	"big"
	"sort"
	"fmt"
	"os"
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

func DotOnLine(p Point, line Line, r int64) bool {
	v1 := NewVector(line[0], line[1])
	v2 := NewVector(line[0], p)
	v3 := VectorProduct(v1, v2)
	l3 := len2(v3)
	l3.Mul(l3, big.NewInt(4))
	r2 := big.NewInt(r)
	r2.Mul(r2, r2)
	r2.Mul(r2, len2(v1))
	return r2.Cmp(l3) >= 0
}

func IntersectLines(l1, l2 Line, scale int64) (p Point, ok bool) {
	panic("IntersectLines is not implemented")
}

func CrossLineWithTriangle(line Line, triangle Triangle, scale int64) (res []Point) {
	fmt.Fprintf(os.Stderr, "CrossLineWithTriangle, line: %v, triangle: %v\n", line, triangle)
	if DotInPlane(line[0], triangle[0], triangle[1], triangle[2], 1) &&
		DotInPlane(line[1], triangle[0], triangle[1], triangle[2], 1) {
		if DotInTriangle(line[0], triangle[0], triangle[1], triangle[2], 1) {
			res = append(res, line[0])
		}
		if DotInTriangle(line[1], triangle[0], triangle[1], triangle[2], 1) {
			res = append(res, line[1])
		}
		if len(res) == 2 {
			return
		}
		before := len(res)
		for i := 0; i < 3; i++ {
			if DotOnLine(triangle[i], line, 1) {
				res = append(res, triangle[i])
			}
		}
		sort.Sort(pointSlice(res))
		res = uniq(res)
		if before < len(res) {
			return
		}
		for i := 0; i < 3; i++ {
			if p, ok := IntersectLines(line, Line{triangle[i], triangle[(i+1)%3]}, scale); ok {
				res = append(res, p)
			}
		}
		sort.Sort(pointSlice(res))
		res = uniq(res)
		return
	}
	return
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

func AllTriangleDots(a, b, c Point, scale, r int64) (res []Point) {
	ga := toGrid(a, scale)
	gb := toGrid(b, scale)
	gc := toGrid(c, scale)

	r = r * scale
	q := []Point{ga, gb, gc}
	var q2 []Point
	m := make(map[uint64]Point)
	m[hash(ga)] = ga
	m[hash(gb)] = gb
	m[hash(gc)] = gc
	for len(q) > 0 {
		q, q2 = q2[0:0], q
		for _, p := range q2 {
			for _, p2 := range adjacent(p) {
				if _, ok := m[hash(p2)]; ok {
					continue
				}

				if !DotInTriangle(scalePoint(p2, scale), a, b, c, r) {
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

func findJ(p1, p2 Point, scale int64) (j uint) {
	for j = 0; j < 31; j++ {
		var r2 int64
		for z := 0; z < 3; z++ {
			diff := int64(p1[z] - p2[z])
			r2 += diff * diff
		}
		//		fmt.Fprintf(os.Stderr, "r2: %d, j: %d, scale: %d\n", r2, j, scale)
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
	//	defer func() {
	//		fmt.Fprintf(os.Stderr, "Less(%d,%d): %v, %v\n", i, j, res, ps)
	//	}()
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
	//	fmt.Fprintf(os.Stderr, "Swap(%d, %d), before: %v\n", i, j, ps)
	ps[i], ps[j] = ps[j], ps[i]
	//	fmt.Fprintf(os.Stderr, "Swap(%d, %d),  after: %v\n", i, j, ps)
}

func uniq(ps []Point) (res []Point) {
	res = ps[:0]
	for i, p := range ps {
		if i > 0 && peq(ps[i-1], p) {
			continue
		}
		res = append(res, ps[i])
	}
	return
}

func scoreDiff(p1, p2 Point) (res int) {
	for i := 0; i < 3; i++ {
		if p1[i] != p2[i] {
			res++
		}
	}
	return
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
	//			fmt.Fprintf(os.Stderr, "AllTriangleDots1, 60, p=%v\n", p)
	p = toGrid(p, scale)
	/*	var last2 Point
		if i1 != 0 {
			if i1 < 0 {
				fmt.Fprintf(os.Stderr, "ogogo! i1 < 0, i1: %d\n", i1)
				panic("aaaa!")
			}
			if scoreDiff(last1, p) > 1 {
				//			fmt.Fprintf(os.Stderr, "So, there is a problem; i1: %d, j1: %d\n", i1, j1)
				var delta uint
				for j1+delta <= MaxJ {
					delta++
					last2 = AddDot(a, b, c, scale, vol, i0, i1*(1<<delta)-1, j0, j1+delta, last1)
					//				fmt.Fprintf(os.Stderr, "last1: %v, p: %v, last2: %v\n", last1, p, last2)
					if !peq(p, last2) {
						break
					}
				}
				//			if j1 > MaxJ {
				//				fmt.Fprintf(os.Stderr, "%d = j1 > MaxJ = %d\n", j1, MaxJ)
				//			}
			}
			//		if peq(p, last2) {
			//			fmt.Fprintf(os.Stderr, "opa! last2 == p\n")
			//		}
		}*/
	vol.Set(int(p[0]), int(p[1]), int(p[2]), color)
	//	if scoreDiff(last1, p) > 1 && scoreDiff(last2, p) > 1 {
	//		fmt.Fprintf(os.Stderr, "Returning bad result. last1: %v, p: %v, last2: %v, i1: %d, j1: %d\n", last1, p, last2, i1, j1)
	//	}
	return p
}

func AllTriangleDots1(a, b, c Point, scale int64, vol VolumeSetter, color uint16) {
	//	fmt.Fprintf(os.Stderr, "AllTriangleDots1, 0, a=%v, b=%v, c=%v\n", a, b, c)
	j0 := findJ(a, c, scale)
	//	fmt.Fprintf(os.Stderr, "AllTriangleDots1, 10, j0=%d\n", j0)
	j1 := findJ(a, b, scale)
	//	fmt.Fprintf(os.Stderr, "AllTriangleDots1, 20, j1=%d\n", j1)
	m := j0
	if m < j1 {
		m = j1
	}
	//	fmt.Fprintf(os.Stderr, "AllTriangleDots1, 30, m=%d\n", m)
	for i0 := 0; i0 <= 1<<j0; i0++ {
		var last1 Point
		for i1 := 0; i0*(1<<(m-j0))+i1*(1<<(m-j1)) <= 1<<m; i1++ {
			last1 = AddDot(a, b, c, scale, vol, int64(i0), int64(i1), j0, j1, last1, color)

		}
		//		fmt.Fprintf(os.Stderr, "AllTriangleDots1, 90\n")
	}
	//	fmt.Fprintf(os.Stderr, "AllTriangleDots1, 94, res: %v\n", res)
	//	fmt.Fprintf(os.Stderr, "AllTriangleDots1, 95, res: %v\n", res)
	//	fmt.Fprintf(os.Stderr, "AllTriangleDots1, 100, res: %v\n", res)
}

func checkAlphaInd(num, den int64, a, b, p, q *Point, ind int) bool {
	if den == 0 {
		return false
	}
	if den < 0 {
		num = -num
		den = -den
	}
	if num < 0 || num > den {
		// 0 <= \alpha <= 1
		return false
	}
	left := a[ind]*den + num*(b[ind]-a[ind])
	if left < p[ind]*den {
		return false
	}
	if left > q[ind]*den {
		return false
	}
	return true
}

func checkAlpha(num, den int64, a, b, p, q *Point) bool {
	for i := 0; i < 3; i++ {
		if !checkAlphaInd(num, den, a, b, p, q, i) {
			return false
		}
	}
	return true
}

func getAlphaPoint(num, den int64, a, b *Point) (res Point) {
	for i := 0; i < 3; i++ {
		// This is not the best thing to do, because we divide on the calculated value.
		// This can lead to an unpredictable behavior, but we say "fine" for now.
		res[i] = a[i] + (num*(b[i]-a[i]))/den
	}
	return
}

// Assumtions: line and cube are non-point, cube[1][i] >= cube[0][i], i \in [0,2]
func ClipLine(line Line, cube Cube, scale int64) (res Line, ok bool) {
	p := scalePoint(cube[0], scale)
	q := scalePoint(cube[1], scale)
	a := scalePoint(line[0], scale)
	b := scalePoint(line[1], scale)
	var num []int64
	var den []int64
	if peq(p, q) || peq(a, b) {
		return
	}

	try := func(n, d int64) {
		if d < 0 {
			n = -n
			d = -d
		}
		if checkAlpha(n, d, &a, &b, &p, &q) {
			num = append(num, n)
			den = append(den, d)
		}
	}
	for i := 0; i < 3; i++ {
		try(p[i]-a[i], b[i]-a[i])
		try(q[i]-a[i], b[i]-a[i])
	}
	try(0, 1)
	try(1, 1)
	//	fmt.Printf("num: %v, den: %v\n", num, den)

	ind := -1
	for i := 1; i < len(num); i++ {
		if num[i]*den[0] != num[0]*den[i] {
			ind = i
			break
		}
	}
	if ind == -1 {
		return
	}
	num[1] = num[ind]
	den[1] = den[ind]
	num = num[:2]
	den = den[:2]
	//	fmt.Printf("ready, num: %v, den: %v\n", num, den)
	if num[1]*den[0] < num[0]*den[1] {
		num[0], num[1] = num[1], num[0]
		den[0], den[1] = den[1], den[0]
	}
	res[0] = toGrid(getAlphaPoint(num[0], den[0], &a, &b), scale)
	res[1] = toGrid(getAlphaPoint(num[1], den[1], &a, &b), scale)

	return res, true
}

var edges = [][6]int64{
	{0, 0, 0, 1, 0, 0},
	{0, 0, 0, 0, 1, 0},
	{0, 0, 0, 0, 0, 1},
	{0, 1, 1, 1, 1, 1},
	{0, 1, 1, 0, 0, 1},
	{0, 1, 1, 0, 1, 0},
	{1, 1, 0, 0, 1, 0},
	{1, 1, 0, 1, 0, 0},
	{1, 1, 0, 1, 1, 1},
	{1, 0, 1, 0, 0, 1},
	{1, 0, 1, 1, 1, 1},
	{1, 0, 1, 1, 0, 0},
}

func ClipTriangle(triangle Triangle, cube Cube, scale int64) (res []Point, ok bool) {
	for i := 0; i < 3; i++ {
		line := Line{triangle[i], triangle[(i+1)%3]}
		if p, ok := ClipLine(line, cube, scale); ok {
			res = append(res, p[0], p[1])
		}
	}
	// It's also possible that triangle crosses one of cube's edges.
	// In this case, the crossing point becomes one of the vertices of the
	// clipped polygon.
	for _, e := range edges {
		var edge Line
		for i := 0; i < 3; i++ {
			edge[0][i] = (1-e[i])*cube[0][i] + e[i]*cube[1][i]
			edge[1][i] = (1-e[i+3])*cube[0][i] + e[i+3]*cube[1][i]
		}
		// TODO: we need to add cross points to the proper place 
		res = append(res, CrossLineWithTriangle(edge, triangle, scale)...)
	}
	res = uniq(res)
	if len(res) > 1 && peq(res[0], res[len(res)-1]) {
		res = res[:len(res)-1]
	}
	return res, len(res) > 2
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
