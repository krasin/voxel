package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/krasin/g3"
	"github.com/krasin/stl"
	//	"github.com/krasin/voxel/nptl"
	"github.com/krasin/voxel/raster"
	"github.com/krasin/voxel/surface"
	"github.com/krasin/voxel/timing"
	"github.com/krasin/voxel/triangle"
	"github.com/krasin/voxel/volume"
)

const (
	VoxelSide      = 512
	MeshMultiplier = 2048
)

var (
	// For neighbours4
	n6dx = []int{1, 0, -1, 0, 0, 0}
	n6dy = []int{0, 1, 0, -1, 0, 0}
	n6dz = []int{0, 0, 0, 0, 1, -1}

	n4dx = n6dx[0:4]
	n4dy = n6dy[0:4]
)

func Index(vol volume.Space16, x, y, z int) int {
	return x*vol.N()*vol.N() + y*vol.N() + z
}

func Coord(vol volume.Space16, index int) (x, y, z int) {
	z = index % vol.N()
	index /= vol.N()
	y = index % vol.N()
	index /= vol.N()
	x = index
	return
}

func Optimize(vol volume.Space16, n int) {
	var q, q2 []int
	vol.SetAllFilled(1, math.MaxUint16-3)
	vol.MapBoundary(func(node g3.Node) {
		if node[2] > 8 {
			vol.Set16(node, 1)
			q = append(q, Index(vol, node[0], node[1], node[2]))
			return
		}
	})

	for len(q) > 0 {
		fmt.Fprintf(os.Stderr, "len(q): %d\n", len(q))
		q, q2 = q2[:0], q
		for _, index := range q2 {
			x, y, z := Coord(vol, index)
			v := vol.Get16(g3.Node{x, y, z})
			if v == 0 {
				panic(fmt.Sprintf("x: %d, y: %d, z: %d, v == 0", x, y, z))
			}
			for k := 0; k < 6; k++ {
				x1 := x + n6dx[k]
				y1 := y + n6dy[k]
				z1 := z + n6dz[k]
				v1 := vol.Get16(g3.Node{x1, y1, z1})
				if v1 > v+1 && int(v)+1 <= n {
					vol.Set16(g3.Node{x1, y1, z1}, v+1)
					q = append(q, Index(vol, x1, y1, z1))
				}
			}
		}
	}
	vol.SetAllFilled(uint16(n+1), 0)

	return
}

type Location16 [2]int16

var sSourcePoint = [3]surface.Vector{
	{0.35, 0.35, 0.35},
	{0.35, 0.65, 0.35},
	{0.65, 0.35, 0.65},
}

func sampleField(fX, fY, fZ float64) float64 {
	if math.Abs(fX) < 0.1 || math.Abs(fY) < 0.1 || math.Abs(fZ) < 0.1 ||
		math.Abs(fX) > 0.9 || math.Abs(fY) > 0.9 || math.Abs(fZ) > 0.9 {
		return 0
	}
	fResult := 0.0
	fDx := fX - sSourcePoint[0].X
	fDy := fY - sSourcePoint[0].Y
	fResult += 0.5 / (fDx*fDx + fDy*fDy)

	fDx = fX - sSourcePoint[1].X
	fDz := fZ - sSourcePoint[1].Z
	fResult += 0.75 / (fDx*fDx + fDz*fDz)

	fDy = fY - sSourcePoint[2].Y
	fDz = fZ - sSourcePoint[2].Z
	fResult += 1.0 / (fDy*fDy + fDz*fDz)

	return fResult
}

func sampleField2(fX, fY, fZ float64) float64 {
	var fResult float64
	var fDx, fDy, fDz float64
	fDx = fX - sSourcePoint[0].X
	fDy = fY - sSourcePoint[0].Y
	fDz = fZ - sSourcePoint[0].Z
	fResult += 0.5 / (fDx*fDx + fDy*fDy + fDz*fDz)

	fDx = fX - sSourcePoint[1].X
	fDy = fY - sSourcePoint[1].Y
	fDz = fZ - sSourcePoint[1].Z
	fResult += 1.0 / (fDx*fDx + fDy*fDy + fDz*fDz)

	fDx = fX - sSourcePoint[2].X
	fDy = fY - sSourcePoint[2].Y
	fDz = fZ - sSourcePoint[2].Z
	fResult += 1.5 / (fDx*fDx + fDy*fDy + fDz*fDz)

	return fResult
}

func NewVolumeField(vol volume.Space16) g3.ScalarField {
	return func(p g3.Point) float64 {
		for _, v := range p {
			if v <= 0 || v >= 1 {
				return 0
			}
		}
		xx := int(p[0] * float64(vol.N()))
		yy := int(p[1] * float64(vol.N()))
		zz := int(p[2] * float64(vol.N()))
		val := vol.Get(g3.Node{xx, yy, zz})
		if val {
			return 100
		}
		return 0
	}
}

type adj struct {
	dx, dy, dz int
	weight     float64
}

var cells = []adj{
	{0, 0, 0, 0.85},
	{1, 0, 0, 0.5},
	{0, 1, 0, 0.5},
	{0, 0, 1, 0.5},
	{-1, 0, 0, 0.5},
	{0, -1, 0, 0.5},
	{0, 0, -1, 0.5},
	{1, 1, 0, 0.25},
	{1, 0, 1, 0.25},
	{0, 1, 1, 0.25},
	{-1, 1, 0, 0.25},
	{-1, 0, 1, 0.25},
	{0, -1, 1, 0.25},
	{1, -1, 0, 0.25},
	{1, 0, -1, 0.25},
	{0, 1, -1, 0.25},
	{-1, -1, 0, 0.25},
	{-1, 0, -1, 0.25},
	{0, -1, -1, 0.25},
}

func NewVolumeField2(vol volume.Space16) g3.ScalarField {
	return func(p g3.Point) float64 {
		for _, v := range p {
			if v <= 0 || v >= 1 {
				return 0
			}
		}

		fx := p[0] * float64(vol.N())
		fy := p[1] * float64(vol.N())
		fz := p[2] * float64(vol.N())

		xx := int(fx)
		yy := int(fy)
		zz := int(fz)

		dx := fx - float64(xx) - 0.5
		dy := fy - float64(yy) - 0.5
		dz := fz - float64(zz) - 0.5

		//		r02 := dx*dx + dy*dy + dz*dz + 0.1

		var val float64
		for _, cell := range cells {
			var v float64
			if vol.Get(g3.Node{xx + cell.dx, yy + cell.dy, zz + cell.dz}) {
				v = 1
			}
			ddx := dx - float64(cell.dx)
			ddy := dy - float64(cell.dy)
			ddz := dz - float64(cell.dz)
			r2 := ddx*ddx + ddy*ddy + ddz*ddz + 0.1
			val += cell.weight * v / r2
		}
		/*		v0 := float64(0)
				if vol.Get(xx, yy, zz) {
					v0 = 1
				}
				val := 0.85 * v0 / r02*/
		return val
	}
}

func main() {
	timing.StartTiming("total")
	timing.StartTiming("Read STL from Stdin")
	triangles, err := stl.Read(os.Stdin)
	if err != nil {
		log.Fatalf("stl.Read: %v", err)
	}
	timing.StopTiming("Read STL from Stdin")

	timing.StartTiming("STLToMesh")
	mesh := raster.STLToMesh(VoxelSide*MeshMultiplier, triangles)
	timing.StopTiming("STLToMesh")

	timing.StartTiming("MeshVolume")
	volume := triangle.MeshVolume(mesh.Triangle, 1)
	if volume < 0 {
		volume = -volume
	}
	fmt.Fprintf(os.Stderr, "Mesh volume (in mesh units): %d\n", volume)
	timing.StopTiming("MeshVolume")

	timing.StartTiming("Rasterize")
	vol := raster.Rasterize(mesh, VoxelSide)
	timing.StopTiming("Rasterize")

	timing.StartTiming("Optimize")
	Optimize(vol, 22)
	timing.StopTiming("Optimize")

	/*	timing.StartTiming("Write nptl")
		if err = nptl.Write(os.Stdout, vol, mesh.Grid); err != nil {
			log.Fatalf("nptl.Write: %v", err)
		}
		v := vol.Volume()
		fmt.Fprintf(os.Stderr, "Volume is filled by %v%%\n", float64(v)*float64(100)/(float64(vol.N())*float64(vol.N())*float64(vol.N())))
		timing.StopTiming("Write nptl")
	*/

	//	t := surface.MarchingCubes(sampleField2, 256, 48)
	side := mesh.Grid.Side()
	vsize := surface.Vector{side, side, side}
	t := surface.MarchingCubes(NewVolumeField2(vol), 128, 0.8, vsize)
	var f *os.File
	if f, err = os.OpenFile("output.stl", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		log.Fatal(err)
	}
	if err = stl.Write(f, t); err != nil {
		log.Fatalf("stl.Write: %v", err)
	}
	f.Close()

	timing.StopTiming("total")
	timing.PrintTimings(os.Stderr)
}
