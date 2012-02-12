package nptl

import (
	"fmt"
	"io"

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

func WriteNptl(vol Uint16Volume, grid Grid, output io.Writer) (err error) {
	stepX := (grid.P1[0] - grid.P0[0]) / float64(vol.XLen())
	stepY := (grid.P1[1] - grid.P0[1]) / float64(vol.YLen())
	stepZ := (grid.P1[2] - grid.P0[2]) / float64(vol.ZLen())

	vol.MapBoundary(func(x, y, z int) {
		nx, ny, nz := Normal(vol, x, y, z)
		if _, err = fmt.Fprintf(output, "%f %f %f %f %f %f\n",
			grid.P0[0]+float64(x)*stepX,
			grid.P0[1]+float64(y)*stepY,
			grid.P0[2]+float64(z)*stepZ,
			nx, ny, nz); err != nil {
			return
		}
	})
	return
}
