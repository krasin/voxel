package nptl

import (
	"fmt"
	"io"

	"github.com/krasin/g3"
	"github.com/krasin/voxel/raster"
	"github.com/krasin/voxel/volume"
)

func Write(w io.Writer, vol volume.Space16, grid raster.Grid) (err error) {
	stepX := (grid.P1[0] - grid.P0[0]) / float64(vol.N())
	stepY := (grid.P1[1] - grid.P0[1]) / float64(vol.N())
	stepZ := (grid.P1[2] - grid.P0[2]) / float64(vol.N())

	vol.MapBoundary(func(node g3.Node) {
		nv := volume.Normal(vol, node)
		if _, err = fmt.Fprintf(w, "%f %f %f %f %f %f\n",
			grid.P0[0]+float64(node[0])*stepX,
			grid.P0[1]+float64(node[1])*stepY,
			grid.P0[2]+float64(node[2])*stepZ,
			nv[0], nv[1], nv[2]); err != nil {
			return
		}
	})
	return
}
