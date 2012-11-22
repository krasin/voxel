package raster

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/krasin/g3"
	"github.com/krasin/stl"
	"github.com/krasin/voxel/set"
	"github.com/krasin/voxel/timing"
	"github.com/krasin/voxel/triangle"
	"github.com/krasin/voxel/volume"
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
)

func PinkX(n int) color.RGBA {
	if n < 200 {
		return color.RGBA{Pink.R, Pink.G + uint8(n), Pink.B, 255}
	}
	return Pink
}

type Mesh struct {
	g3.Grid
	Triangle []triangle.Triangle
}

func STLToMesh(n int, triangles []stl.Triangle) (m Mesh) {
	if n < 4 {
		panic("n < 4")
	}
	// Find bounds
	min := g3.Point{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	max := g3.Point{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	for _, t := range triangles {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				cur := float64(t.V[i][j])
				if min[j] > cur {
					min[j] = cur
				}
				if max[j] < cur {
					max[j] = cur
				}
			}
		}
	}

	// Now, we need to determine a step H for uniform grid that starts from point min,
	// which will contains all the figure.
	// Additional trick: we want min to be at (1,1,1) and max not far than (n-1,n-1,n-1)
	// to prevent edge effects.
	vsize := max.Sub(min)
	var maxSize float64
	for _, v := range vsize {
		if v > maxSize {
			maxSize = v
		}
	}
	h := maxSize / float64(n-2)

	// Now, we know all mesh parameters: N, H and know that min is at (1,1,1)
	m.N = n
	m.H = h
	m.P0 = min.Sub(g3.Point{h, h, h})
	m.Triangle = make([]triangle.Triangle, len(triangles))
	for i, t := range triangles {
		cur := &m.Triangle[i]
		for i := 0; i < 3; i++ {
			p := g3.Point{float64(t.V[i][0]), float64(t.V[i][1]), float64(t.V[i][2])}
			node := m.Node(p)
			cur[i][0] = int64(node[0])
			cur[i][1] = int64(node[1])
			cur[i][2] = int64(node[2])
		}
	}
	return
}

func Rasterize(m Mesh, n int) volume.Space16 {
	scale := m.N / n
	vol := volume.NewSparseVolume(n)

	timing.StartTiming("Rasterize triangles")
	for index, t := range m.Triangle {
		triangle.AllTriangleDots(t[0], t[1], t[2], int64(scale), vol, uint16(1+(index%10)))
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
				p2[i] = p2[i] + j
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
					p2[i] = p2[i] + j
					color2 := vol.Get16(g3.Node{int(p2[0]), int(p2[1]), int(p2[2])})
					if int(color2) < shift {
						continue
					}
					if color == 0 {
						vol.Set16(g3.Node{int(p[0]), int(p[1]), int(p[2])}, color2)
						color = color2
					} else {
						ds.Join(int(color)-shift, int(color2)-shift)
					}
				}
			}
			if color == 0 {
				vol.Set16(g3.Node{int(p[0]), int(p[1]), int(p[2])}, uint16(shift+ds.Make()))
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
					v := vol.Get16(g3.Node{x, y, z})
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
