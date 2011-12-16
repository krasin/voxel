package main

import (
	"fmt"
	"io"
	"image"
	"image/png"
	"log"
	"math"
	"os"
)

const (
	SizeOfSTLTriangle = 4*3*4 + 2
	VoxelSide         = 512
	MeshMultiplier    = 2048
)

var (
	Black  = image.RGBAColor{0, 0, 0, 255}
	Yellow = image.RGBAColor{255, 255, 0, 255}
	Green  = image.RGBAColor{0, 255, 0, 255}
	colors = []image.RGBAColor{
		{0, 0, 0, 255},
		{255, 255, 0, 255},
		{0, 255, 0, 255},
		{255, 0, 0, 255},
		{0, 0, 255, 255},
		{128, 255, 0, 255},
		{0, 255, 128, 255},
		{255, 128, 0, 255},
		{255, 0, 128, 255},
		{0, 128, 255, 255},
		{128, 0, 255, 255},
		{128, 128, 255, 255},
	}

	// For neighbours4
	n6dx = []int{1, 0, -1, 0, 0, 0}
	n6dy = []int{0, 1, 0, -1, 0, 0}
	n6dz = []int{0, 0, 0, 0, 1, -1}

	n4dx = n6dx[0:4]
	n4dy = n6dy[0:4]
)

type BoolVoxelVolume interface {
	Get(x, y, z int) bool
	XLen() int
	YLen() int
	ZLen() int
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

func Normal(vol BoolVoxelVolume, x, y, z int) (nx, ny, nz float64) {
	var px, py, pz int
	if vol.Get(x-1, y, z) {
		px++
	}
	if vol.Get(x+1, y, z) {
		px--
	}
	if vol.Get(x, y-1, z) {
		py++
	}
	if vol.Get(x, y+1, z) {
		py--
	}
	if vol.Get(x, y, z-1) {
		pz++
	}
	if vol.Get(x, y, z+1) {
		pz--
	}
	if vol.Get(x-1, y-1, z) {
		px++
		py++
	}
	if vol.Get(x+1, y-1, z) {
		px--
		py++
	}
	if vol.Get(x-1, y+1, z) {
		px++
		py--
	}
	if vol.Get(x+1, y+1, z) {
		px--
		py--
	}
	if vol.Get(x, y-1, z-1) {
		py++
		pz++
	}
	if vol.Get(x, y+1, z-1) {
		py--
		pz++
	}
	if vol.Get(x, y-1, z+1) {
		py++
		pz--
	}
	if vol.Get(x, y+1, z+1) {
		py--
		pz--
	}
	if vol.Get(x-1, y, z-1) {
		px++
		pz++
	}
	if vol.Get(x+1, y, z-1) {
		px--
		pz++
	}
	if vol.Get(x-1, y, z+1) {
		px++
		pz--
	}
	if vol.Get(x+1, y, z+1) {
		px--
		pz--
	}
	r2 := px*px + py*py + pz*pz
	if r2 == 0 {
		return 1, 0, 0
	}
	l := math.Sqrt(float64(r2))
	nx = float64(px) / l
	ny = float64(py) / l
	nz = float64(pz) / l
	return
}

type Uint16Volume interface {
	BoolVoxelVolume
	Set(x, y, z int, v uint16)
	GetV(x, y, z int) uint16
}

func Index(vol Uint16Volume, x, y, z int) int {
	return x*vol.YLen()*vol.ZLen() + y*vol.ZLen() + z
}

func Coord(vol Uint16Volume, index int) (x, y, z int) {
	z = index % vol.ZLen()
	index /= vol.ZLen()
	y = index % vol.YLen()
	index /= vol.YLen()
	x = index
	return
}

func Optimize(vol Uint16Volume, n int) {
	var q, q2 []int
	for y := 0; y < vol.YLen(); y++ {
		for z := 0; z < vol.ZLen(); z++ {
			for x := 0; x < vol.XLen(); x++ {
				if IsBoundary(vol, x, y, z) && z > 8 {
					vol.Set(x, y, z, 1)
					q = append(q, Index(vol, x, y, z))
					continue
				}
				if vol.Get(x, y, z) {
					vol.Set(x, y, z, math.MaxUint16-3)
				}
			}
		}
	}
	for len(q) > 0 {
		fmt.Fprintf(os.Stderr, "len(q): %d\n", len(q))
		q, q2 = q2[:0], q
		for _, index := range q2 {
			x, y, z := Coord(vol, index)
			v := vol.GetV(x, y, z)
			if v == 0 {
				panic(fmt.Sprintf("x: %d, y: %d, z: %d, v == 0", x, y, z))
			}
			for k := 0; k < 6; k++ {
				x1 := x + n6dx[k]
				y1 := y + n6dy[k]
				z1 := z + n6dz[k]
				v1 := vol.GetV(x1, y1, z1)
				if v1 > v+1 && int(v)+1 <= n {
					vol.Set(x1, y1, z1, v+1)
					q = append(q, Index(vol, x1, y1, z1))
				}
			}
		}
	}
	for y := 0; y < vol.YLen(); y++ {
		for z := 0; z < vol.ZLen(); z++ {
			for x := 0; x < vol.XLen(); x++ {
				if vol.GetV(x, y, z) > uint16(n) {
					vol.Set(x, y, z, 0)
				}
			}
		}
	}
	return
}

func WriteNptl(vol BoolVoxelVolume, grid Grid, output io.Writer) (err os.Error) {
	v := 0
	stepX := (grid.P1[0] - grid.P0[0]) / float64(vol.XLen())
	stepY := (grid.P1[1] - grid.P0[1]) / float64(vol.YLen())
	stepZ := (grid.P1[2] - grid.P0[2]) / float64(vol.ZLen())
	for y := 0; y < vol.YLen(); y++ {
		for z := 0; z < vol.ZLen(); z++ {
			for x := 0; x < vol.XLen(); x++ {
				if !vol.Get(x, y, z) {
					continue
				}
				v++
				if !IsBoundary(vol, x, y, z) {
					continue
				}
				nx, ny, nz := Normal(vol, x, y, z)
				if _, err = fmt.Fprintf(output, "%f %f %f %f %f %f\n",
					grid.P0[0]+float64(x)*stepX,
					grid.P0[1]+float64(y)*stepY,
					grid.P0[2]+float64(z)*stepZ,
					nx, ny, nz); err != nil {
					return
				}
			}
		}
	}
	fmt.Fprintf(os.Stderr, "Volume is filled by %v%\n", float64(v)*float64(100)/(float64(vol.XLen())*float64(vol.YLen())*float64(vol.ZLen())))
	return
}

type Grid struct {
	P0 [3]float64
	P1 [3]float64
	N  [3]int64
}

type Mesh struct {
	Grid
	Triangle []Triangle
}

func STLToMesh(n int, triangles []STLTriangle) (m Mesh) {
	min := []float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	max := []float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
	for _, t := range triangles {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				cur := t.v[i][j]
				if min[j] > cur {
					min[j] = cur
				}
				if max[j] < cur {
					max[j] = cur
				}
			}
		}
	}
	m.P0 = [3]float64{float64(min[0]) - 1, float64(min[1]) - 1, float64(min[2]) - 1}
	m.P1 = [3]float64{float64(max[0]) + 1, float64(max[1]) + 1, float64(max[2]) + 1}
	m.N = [3]int64{int64(n), int64(n), int64(n)}
	m.Triangle = make([]Triangle, len(triangles))
	for i, t := range triangles {
		cur := &m.Triangle[i]
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				val := float64(t.v[i][j])
				cur[i][j] = int64((val - m.P0[j]) * float64(m.N[j]) / (m.P1[j] - m.P0[j]))
			}
		}
	}
	fmt.Fprintf(os.Stderr, "min: %v, max: %v\n", min, max)
	fmt.Fprintf(os.Stderr, "mesh.Grid: %v, mesh.Triangle: %v\n", m.Grid, m.Triangle[:10])
	return
}

type VolumeSetter interface {
	Set(x, y, z int, val uint16)
}

type Location16 [2]int16

func Rasterize(m Mesh, n int) Uint16Volume {
	scale := m.N[0] / int64(n)
	vol := NewOctree(n)
	// Rasterize edges
	for index, t := range m.Triangle {
		AllTriangleDots(t[0], t[1], t[2], scale, vol, uint16(1+(index%10)))
	}
	fmt.Fprintf(os.Stderr, "Triangle rasterization complete\n")

	var cnt int
	bmp := image.NewRGBA(n, n)
	var q, q2 []Location16
	for z := 1; z < n; z++ {
		cnt = 0
		q = q[:0]
		q2 = q2[:0]
		for x := 0; x < n; x++ {
			for y := 0; y < n; y++ {
				if !vol.Get(x, y, z) {
					if x > 0 && y > 0 && z > 0 && x < n-1 && y < n-1 && z < n-1 {
						vol.Set(x, y, z, 11)
					} else {
						q = append(q, Location16{int16(x), int16(y)})
					}
				}
			}
		}
		for len(q) > 0 {
			q, q2 = q2[:0], q
			for _, cur := range q2 {
				for i := 0; i < 4; i++ {
					x1 := int(cur[0]) + n4dx[i]
					y1 := int(cur[1]) + n4dy[i]
					if x1 <= 0 || x1 >= n-1 || y1 <= 0 || y1 >= n-1 {
						continue
					}
					if vol.GetV(x1, y1, z) == 11 {
						vol.Set(x1, y1, z, 0)
						q = append(q, Location16{int16(x1), int16(y1)})
					}
				}
			}
		}
		if z%10 == 0 {
			for x := 0; x < n; x++ {
				for y := 0; y < n; y++ {
					v := vol.GetV(x, y, z)
					if v <= 11 {
						bmp.Set(x, y, colors[v])
					} else {
						bmp.Set(x, y, Black)
					}
				}
			}

			f, _ := os.OpenFile(fmt.Sprintf("zban-%03d.png", z), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
			png.Encode(f, bmp)
			f.Close()
		}
	}
	fmt.Fprintf(os.Stderr, "Last layer: %d dots\n", cnt)
	return vol
}

func main() {
	triangles, err := ReadSTL(os.Stdin)
	if err != nil {
		log.Fatalf("ReadSTL: %v", err)
	}
	mesh := STLToMesh(VoxelSide*MeshMultiplier, triangles)

	volume := MeshVolume(mesh.Triangle, 1)
	if volume < 0 {
		volume = -volume
	}
	fmt.Fprintf(os.Stderr, "Mesh volume (in mesh units): %d\n", volume)
	fmt.Fprintf(os.Stderr, "Mesh volume (original units): %f\n", float64(volume)/float64(mesh.N[0]*mesh.N[1]*mesh.N[2])*(mesh.P1[0]-mesh.P0[0])*(mesh.P1[1]-mesh.P0[1])*(mesh.P1[2]-mesh.P0[2]))

	vol := Rasterize(mesh, VoxelSide)

	Optimize(vol, 22)
	if err = WriteNptl(vol, mesh.Grid, os.Stdout); err != nil {
		log.Fatalf("WriteNptl: %v", err)
	}
}
