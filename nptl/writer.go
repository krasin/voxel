package nptl

import (
	"fmt"
	"io"

	"github.com/krasin/g3"
	"github.com/krasin/voxel/raster"
	"github.com/krasin/voxel/volume"
)

type Uint16Volume interface {
	volume.BoolVoxelVolume
	Set(x, y, z int, v uint16)
	GetV(x, y, z int) uint16
	SetAllFilled(threshold, val uint16)
	MapBoundary(f func(x, y, z int))
	Volume() int64
}

func Write(w io.Writer, vol Uint16Volume, grid raster.Grid) (err error) {
	stepX := (grid.P1[0] - grid.P0[0]) / float64(vol.N())
	stepY := (grid.P1[1] - grid.P0[1]) / float64(vol.N())
	stepZ := (grid.P1[2] - grid.P0[2]) / float64(vol.N())

	vol.MapBoundary(func(x, y, z int) {
		nv := volume.Normal(vol, g3.Node{x, y, z})
		if _, err = fmt.Fprintf(w, "%f %f %f %f %f %f\n",
			grid.P0[0]+float64(x)*stepX,
			grid.P0[1]+float64(y)*stepY,
			grid.P0[2]+float64(z)*stepZ,
			nv[0], nv[1], nv[2]); err != nil {
			return
		}
	})
	return
}
