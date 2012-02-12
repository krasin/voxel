package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"os"

	"github.com/krasin/stl"
	"github.com/krasin/voxel/nptl"
	"github.com/krasin/voxel/set"
	"github.com/krasin/voxel/timing"
	"github.com/krasin/voxel/triangle"
	"github.com/krasin/voxel/volume"
)

const (
	VoxelSide      = 512
	MeshMultiplier = 2048
)

var (
	Black  = color.RGBA{0, 0, 0, 255}
	Yellow = color.RGBA{255, 255, 0, 255}
	Green  = color.RGBA{0, 255, 0, 255}
	Pink   = color.RGBA{255, 20, 147, 255}
	colors = []color.RGBA{
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
		{255, 255, 128, 255},
		{128, 255, 128, 255},
		{128, 128, 255, 255},
	}

	// For neighbours4
	n6dx = []int{1, 0, -1, 0, 0, 0}
	n6dy = []int{0, 1, 0, -1, 0, 0}
	n6dz = []int{0, 0, 0, 0, 1, -1}

	n4dx = n6dx[0:4]
	n4dy = n6dy[0:4]
)

func PinkX(n int) color.RGBA {
	if n < 200 {
		return color.RGBA{Pink.R, Pink.G + uint8(n), Pink.B, 255}
	}
	return Pink
}

type Uint16Volume interface {
	volume.BoolVoxelVolume
	Set(x, y, z int, v uint16)
	GetV(x, y, z int) uint16
	SetAllFilled(threshold, val uint16)
	MapBoundary(f func(x, y, z int))
	Volume() int64
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
	vol.SetAllFilled(1, math.MaxUint16-3)
	vol.MapBoundary(func(x, y, z int) {
		if z > 8 {
			vol.Set(x, y, z, 1)
			q = append(q, Index(vol, x, y, z))
			return
		}
	})

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
	vol.SetAllFilled(uint16(n+1), 0)

	return
}

type Grid struct {
	P0 [3]float64
	P1 [3]float64
	N  [3]int64
}

type Mesh struct {
	Grid
	Triangle []triangle.Triangle
}

func STLToMesh(n int, triangles []stl.STLTriangle) (m Mesh) {
	min := []float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	max := []float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
	for _, t := range triangles {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				cur := t.V[i][j]
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
	m.Triangle = make([]triangle.Triangle, len(triangles))
	for i, t := range triangles {
		cur := &m.Triangle[i]
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				val := float64(t.V[i][j])
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
	vol := volume.NewSparseVolume(n)

	timing.StartTiming("Rasterize triangles")
	for index, t := range m.Triangle {
		triangle.AllTriangleDots(t[0], t[1], t[2], scale, vol, uint16(1+(index%10)))
	}
	fmt.Fprintf(os.Stderr, "Triangle rasterization complete\n")
	timing.StopTiming("Rasterize triangles")

	timing.StartTiming("Rasterize cubes")
	ds := set.NewDisjoinSet()
	// Reserve color for outer space
	ds.Make()

	shift := 11

	// Let's color cubes.
	for k, cube := range vol.Cubes {
		// Skip cubes with leaf voxels
		if cube != nil {
			continue
		}
		p := volume.K2cube(k)

		// If this is a cube at the edge of the space, it's a part of outer space.
		if p[0] == 0 || p[1] == 0 || p[2] == 0 ||
			int(p[0]) == (1<<uint(vol.LK))-1 || int(p[1]) == (1<<uint(vol.LK))-1 || int(p[2]) == (1<<uint(vol.LK))-1 {
			vol.Colors[k] = uint16(shift + ds.Find(0))
			continue
		}

		// Look if any neighbour has already color assigned
		for i := 0; i < 3; i++ {
			for j := -1; j <= 1; j += 2 {
				p2 := p
				p2[i] = uint16(int(p2[i]) + j)
				k2 := volume.Cube2k(p2)
				if k2 >= len(vol.Colors) {
					panic(fmt.Sprintf("k2: %d, len(vol.Colors): %d, len(vol.Cubes): %d, p: %v, p2: %v, k: %d", k2, len(vol.Colors), len(vol.Cubes), p, p2, k))
				}
				if vol.Colors[k2] == 0 {
					continue
				}
				if vol.Colors[k] == 0 {
					vol.Colors[k] = vol.Colors[k2]
				} else {
					ds.Join(int(vol.Colors[k])-shift, int(vol.Colors[k2])-shift)
				}
			}
		}

		// If there's no colored neighbour, introduce a new color.
		if vol.Colors[k] == 0 {
			vol.Colors[k] = uint16(shift + ds.Make())
		}
	}
	timing.StopTiming("Rasterize cubes")
	timing.StartTiming("Rasterize leaf voxels")

	// Now, we need to go through cubes which have leaf voxels
	for k, cube := range vol.Cubes {
		if cube == nil {
			continue
		}
		for h, val := range cube {
			if val != 0 {
				continue
			}
			p := volume.Kh2point(k, h)
			color := val
			// Look for neighbours of this leaf voxel
			for i := 0; i < 3; i++ {
				for j := -1; j <= 1; j += 2 {
					p2 := p
					p2[i] = uint16(int(p2[i]) + j)
					color2 := vol.GetV(int(p2[0]), int(p2[1]), int(p2[2]))
					if int(color2) < shift {
						continue
					}
					if color == 0 {
						vol.Set(int(p[0]), int(p[1]), int(p[2]), color2)
						color = color2
					} else {
						ds.Join(int(color)-shift, int(color2)-shift)
					}
				}
			}
			if color == 0 {
				vol.Set(int(p[0]), int(p[1]), int(p[2]), uint16(shift+ds.Make()))
			}
		}
	}
	timing.StopTiming("Rasterize leaf voxels")

	timing.StartTiming("Rasterize.CanonicalizeColors")
	// Canonicalize colors
	canonicalZero := uint16(shift + ds.Find(0))
	for k, cube := range vol.Cubes {
		if cube == nil {
			vol.Colors[k] = uint16(shift + ds.Find(int(vol.Colors[k])-shift))
			if vol.Colors[k] == canonicalZero {
				vol.Colors[k] = 0
			}
			continue
		}
		for h := 0; h < len(cube); h++ {
			if int(cube[h]) < shift {
				continue
			}
			cube[h] = uint16(shift + ds.Find(int(cube[h])-shift))
			if cube[h] == canonicalZero {
				cube[h] = 0
			}
		}
	}
	timing.StopTiming("Rasterize.CanonicalizeColors")

	timing.StartTiming("Rasterize.DrawSlices")
	bmp := image.NewRGBA(image.Rect(0, 0, n, n))
	for z := 1; z < n; z++ {
		if z%10 == 0 {
			for x := 0; x < n; x++ {
				for y := 0; y < n; y++ {
					v := vol.GetV(x, y, z)
					if int(v) < len(colors) {
						bmp.Set(x, y, colors[v])
					} else {
						bmp.Set(x, y, PinkX(int(v)))
					}
				}
			}

			f, _ := os.OpenFile(fmt.Sprintf("zban-%03d.png", z), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
			png.Encode(f, bmp)
			f.Close()
		}
	}
	timing.StopTiming("Rasterize.DrawSlices")
	fmt.Fprintf(os.Stderr, "Rasterize complete\n")
	return vol
}

func main() {
	timing.StartTiming("total")
	timing.StartTiming("Read STL from Stdin")
	triangles, err := stl.ReadSTL(os.Stdin)
	if err != nil {
		log.Fatalf("ReadSTL: %v", err)
	}
	timing.StopTiming("Read STL from Stdin")

	timing.StartTiming("STLToMesh")
	mesh := STLToMesh(VoxelSide*MeshMultiplier, triangles)
	timing.StopTiming("STLToMesh")

	timing.StartTiming("MeshVolume")
	volume := triangle.MeshVolume(mesh.Triangle, 1)
	if volume < 0 {
		volume = -volume
	}
	fmt.Fprintf(os.Stderr, "Mesh volume (in mesh units): %d\n", volume)
	fmt.Fprintf(os.Stderr, "Mesh volume (original units): %f\n", float64(volume)/float64(mesh.N[0]*mesh.N[1]*mesh.N[2])*(mesh.P1[0]-mesh.P0[0])*(mesh.P1[1]-mesh.P0[1])*(mesh.P1[2]-mesh.P0[2]))
	timing.StopTiming("MeshVolume")

	timing.StartTiming("Rasterize")
	vol := Rasterize(mesh, VoxelSide)
	timing.StopTiming("Rasterize")

	timing.StartTiming("Optimize")
	Optimize(vol, 22)
	timing.StopTiming("Optimize")

	timing.StartTiming("WriteNptl")
	if err = WriteNptl(vol, mesh.Grid, os.Stdout); err != nil {
		log.Fatalf("WriteNptl: %v", err)
	}
	v := vol.Volume()
	fmt.Fprintf(os.Stderr, "Volume is filled by %v%%\n", float64(v)*float64(100)/(float64(vol.XLen())*float64(vol.YLen())*float64(vol.ZLen())))
	timing.StopTiming("WriteNptl")

	timing.StopTiming("total")
	timing.PrintTimings(os.Stderr)
}
