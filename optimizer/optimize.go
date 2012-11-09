package main

import (
	"fmt"
	"log"
	"math"
	"os"

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

type VolumeSetter interface {
	Set(x, y, z int, val uint16)
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
	fmt.Fprintf(os.Stderr, "Mesh volume (original units): %f\n", float64(volume)/float64(mesh.N[0]*mesh.N[1]*mesh.N[2])*(mesh.P1[0]-mesh.P0[0])*(mesh.P1[1]-mesh.P0[1])*(mesh.P1[2]-mesh.P0[2]))
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
		fmt.Fprintf(os.Stderr, "Volume is filled by %v%%\n", float64(v)*float64(100)/(float64(vol.XLen())*float64(vol.YLen())*float64(vol.ZLen())))
		timing.StopTiming("Write nptl")
	*/

	t := surface.MarchingCubes(sampleField2, 256, 48)
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
