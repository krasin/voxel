package main

import (
	"big"
	"sort"
)

type Point [3]int64
type Vector [3]int64
type Triangle [3]Point
type Line [2]Point
type Cube [2]Point

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
			return j + 1
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

func AllTriangleDots1(a, b, c Point, scale, r int64) (res []Point) {
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
	cur0 := -1
	for i0 := 0; i0 <= 1<<j0; i0++ {
		ind0 := cur0
		cur0 = len(res)
		var last1 Point
		for i1 := 0; i0*(1<<(m-j0))+i1*(1<<(m-j1)) <= 1<<m; i1++ {
			if ind0 >= 0 && i1 > 0 {
				ind0++
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
			if ind0 >= 0 && ind0 < cur0 && peq(res[ind0], p) {
				continue
			}
			if i1 > 0 && peq(last1, p) {
				continue
			}
			res = append(res, p)
			last1 = p
		}
		//		fmt.Fprintf(os.Stderr, "AllTriangleDots1, 90\n")
	}
	//	fmt.Fprintf(os.Stderr, "AllTriangleDots1, 94, res: %v\n", res)
	sort.Sort(pointSlice(res))
	//	fmt.Fprintf(os.Stderr, "AllTriangleDots1, 95, res: %v\n", res)
	res = uniq(res)
	//	fmt.Fprintf(os.Stderr, "AllTriangleDots1, 100, res: %v\n", res)
	return
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

func ClipTriangle(triangle Triangle, cube Cube, scale int64) (res []Point, ok bool) {
	for i := 0; i < 3; i++ {
		line := Line{triangle[i], triangle[(i+1)%3]}
		if p, ok := ClipLine(line, cube, scale); ok {
			res = append(res, p[0], p[1])
		}
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
