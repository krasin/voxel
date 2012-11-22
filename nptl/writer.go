package nptl

import (
	"fmt"
	"io"

	"github.com/krasin/g3"
	"github.com/krasin/voxel/volume"
)

func Write(w io.Writer, vol volume.Space16, grid g3.Grid) (err error) {
	vol.MapBoundary(func(node g3.Node) {
		nv := volume.Normal(vol, node)
		cur := grid.At(node)
		if _, err = fmt.Fprintf(w, "%f %f %f %f %f %f\n",
			cur[0], cur[1], cur[2], nv[0], nv[1], nv[2]); err != nil {
			return
		}
	})
	return
}
