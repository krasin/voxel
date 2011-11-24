package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"unsafe"
)

const (
	TAG_END        = 0
	TAG_BYTE       = 1
	TAG_SHORT      = 2
	TAG_INT        = 3
	TAG_LONG       = 4
	TAG_FLOAT      = 5
	TAG_DOUBLE     = 6
	TAG_BYTE_ARRAY = 7
	TAG_STRING     = 8
	TAG_LIST       = 9
	TAG_COMPOUND   = 10

	SizeOfSTLTriangle = 4*3*4 + 2
	VoxelSide         = 1024
	MeshMultiplier    = 1024
)

type NBTReader struct {
	r *bufio.Reader
}

func NewNBTReader(r io.Reader) (nr *NBTReader, err os.Error) {
	var rd io.Reader
	if rd, err = gzip.NewReader(r); err != nil {
		return
	}
	return &NBTReader{r: bufio.NewReader(rd)}, nil
}

type Tag int

func (r *NBTReader) ReadString() (str string, err os.Error) {
	var l int
	if l, err = r.ReadShort(); err != nil {
		return
	}
	data := make([]byte, l)
	if _, err = io.ReadFull(r.r, data); err != nil {
		return
	}
	return string(data), nil
}

func (r *NBTReader) ReadShort() (val int, err os.Error) {
	buf := [2]byte{}
	if _, err = io.ReadFull(r.r, buf[:]); err != nil {
		return
	}
	val = int(buf[1]) + (int(buf[0]) << 8) // Big Endian
	return
}

func (r *NBTReader) ReadInt() (val int, err os.Error) {
	buf := [4]byte{}
	if _, err = io.ReadFull(r.r, buf[:]); err != nil {
		return
	}
	for i := 0; i < 4; i++ {
		val <<= 8
		val += int(buf[i])
	}
	return
}

func (r *NBTReader) ReadTagTyp() (typ byte, err os.Error) {
	typ, err = r.r.ReadByte()
	return
}

func (r *NBTReader) ReadTagName() (typ byte, name string, err os.Error) {
	if typ, err = r.r.ReadByte(); err != nil {
		return
	}
	if typ == TAG_END {
		return
	}
	name, err = r.ReadString()
	return
}

func (r *NBTReader) ReadByteArray() (data []byte, err os.Error) {
	var l int
	if l, err = r.ReadInt(); err != nil {
		return
	}
	data = make([]byte, l)
	_, err = io.ReadFull(r.r, data)
	return
}

type SchematicReader struct {
	r *NBTReader
}

type BoolVoxelVolume interface {
	Get(x, y, z int) bool
	XLen() int
	YLen() int
	ZLen() int
}

func NewSchematicReader(r io.Reader) (sr *SchematicReader, err os.Error) {
	var nr *NBTReader
	if nr, err = NewNBTReader(r); err != nil {
		return
	}
	return &SchematicReader{r: nr}, nil
}

type Entity struct {
	Id string
}

type Schematic struct {
	Width     int
	Length    int
	Height    int
	WEOffsetX int
	WEOffsetY int
	WEOffsetZ int
	Materials string
	Blocks    []byte
	Data      []byte
	Entities  []Entity
}

func (r *SchematicReader) ReadEntity() (entity Entity, err os.Error) {
	for {
		var typ byte
		var name string
		if typ, name, err = r.r.ReadTagName(); err != nil {
			return
		}
		if typ == TAG_END {
			break
		}
		switch name {
		default:
			err = fmt.Errorf("Unknown entity field: %s", name)
			return
		}
	}
	return
}

func (r *SchematicReader) ReadEntities() (entities []Entity, err os.Error) {
	for {
		var typ byte
		if typ, err = r.r.ReadTagTyp(); err != nil {
			return
		}
		if typ == TAG_END {
			break
		}
		if typ == TAG_COMPOUND {
		}
		var entity Entity
		if entity, err = r.ReadEntity(); err != nil {
			return
		}
		entities = append(entities, entity)
	}
	return
}

func (r *SchematicReader) Parse() (s *Schematic, err os.Error) {
	var typ byte
	var name string
	if typ, name, err = r.r.ReadTagName(); err != nil {
		return
	}
	if typ != TAG_COMPOUND {
		return nil, fmt.Errorf("Top level tag must be compound. Got: %d", typ)
	}
	if name != "Schematic" {
		return nil, fmt.Errorf("Unexpected tag name: %s, want: Schematic", name)
	}
	s = new(Schematic)
	for {
		if typ, name, err = r.r.ReadTagName(); err != nil {
			return
		}
		if typ == TAG_END {
			break
		}
		switch name {
		case "Width":
			s.Width, err = r.r.ReadShort()
		case "Length":
			s.Length, err = r.r.ReadShort()
		case "Height":
			s.Height, err = r.r.ReadShort()
		case "Materials":
			s.Materials, err = r.r.ReadString()
		case "Blocks":
			s.Blocks, err = r.r.ReadByteArray()
		case "Data":
			s.Data, err = r.r.ReadByteArray()
		case "WEOffsetX":
			s.WEOffsetX, err = r.r.ReadInt()
		case "WEOffsetY":
			s.WEOffsetY, err = r.r.ReadInt()
		case "WEOffsetZ":
			s.WEOffsetZ, err = r.r.ReadInt()
		case "Entities":
			s.Entities, err = r.ReadEntities()
		default:
			return nil, fmt.Errorf("Unexpected tag: %d, name: %s\n", typ, name)
		}
		if err != nil {
			return
		}
	}
	if s.Materials != "Alpha" {
		return nil, fmt.Errorf("Materials must have 'Alpha' value, got: '%s'", s.Materials)
	}
	return
}

func (s *Schematic) XLen() int {
	return s.Width
}

func (s *Schematic) YLen() int {
	return s.Height
}

func (s *Schematic) ZLen() int {
	return s.Length
}

func (s *Schematic) Get(x, y, z int) bool {
	if x < 0 || y < 0 || z < 0 || x >= s.XLen() || y >= s.YLen() || z >= s.ZLen() {
		return false
	}
	index := y*s.XLen()*s.ZLen() + z*s.XLen() + x
	return s.Blocks[index] != 0
}

func ReadSchematic(input io.Reader) (vol BoolVoxelVolume, err os.Error) {
	var r *SchematicReader
	if r, err = NewSchematicReader(input); err != nil {
		return
	}
	vol, err = r.Parse()
	return
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

type ArrayVolume struct {
	a                []uint32
	xlen, ylen, zlen int
}

func NewArrayVolume(xlen, ylen, zlen int) *ArrayVolume {
	l := xlen * ylen * zlen
	return &ArrayVolume{
		a:    make([]uint32, l),
		xlen: xlen,
		ylen: ylen,
		zlen: zlen,
	}
}

func (v *ArrayVolume) XLen() int {
	return v.xlen
}

func (v *ArrayVolume) YLen() int {
	return v.ylen
}

func (v *ArrayVolume) ZLen() int {
	return v.zlen
}

func (v *ArrayVolume) Index(x, y, z int) int {
	return x*v.ylen*v.zlen + y*v.zlen + z
}

func (v *ArrayVolume) Coord(index int) (x, y, z int) {
	z = index % v.zlen
	index /= v.zlen
	y = index % v.ylen
	index /= v.ylen
	x = index
	return
}

func (v *ArrayVolume) Get(x, y, z int) bool {
	if x < 0 || y < 0 || z < 0 || x >= v.XLen() || y >= v.YLen() || z >= v.ZLen() {
		return false
	}
	return v.a[v.Index(x, y, z)] != 0
}

func (v *ArrayVolume) Set(x, y, z int, val uint32) {
	v.a[v.Index(x, y, z)] = val
}

func (v *ArrayVolume) GetV(x, y, z int) uint32 {
	return v.a[v.Index(x, y, z)]
}

func Optimize(input BoolVoxelVolume, n2 int) BoolVoxelVolume {
	res := NewArrayVolume(input.XLen(), input.YLen(), input.ZLen())
	var q, q2 []int
	for y := 0; y < input.YLen(); y++ {
		for z := 0; z < input.ZLen(); z++ {
			for x := 0; x < input.XLen(); x++ {
				if !input.Get(x, y, z) {
					res.Set(x, y, z, 0)
					continue
				}
				if IsBoundary(input, x, y, z) && z > 0 {
					res.Set(x, y, z, 1)
					q = append(q, res.Index(x, y, z))
					continue
				}
				res.Set(x, y, z, math.MaxUint32)
			}
		}
	}
	for len(q) > 0 {
		fmt.Fprintf(os.Stderr, "len(q): %d\n", len(q))
		q, q2 = q2[:0], q
		for _, index := range q2 {
			x, y, z := res.Coord(index)
			v := res.GetV(x, y, z)
			for dx := -1; dx <= 1; dx++ {
				x1 := x + dx
				for dy := -1; dy <= 1; dy++ {
					y1 := y + dy
					for dz := -1; dz <= 1; dz++ {
						z1 := z + dz
						if !res.Get(x1, y1, z1) || dx == 0 && dy == 0 && dz == 0 {
							continue
						}
						r2 := uint32(dx*dx + dy*dy + dz*dz)
						v1 := res.GetV(x1, y1, z1)
						if v1 > v+r2 {
							res.Set(x1, y1, z1, v+r2)
							q = append(q, res.Index(x1, y1, z1))
						}
					}
				}
			}
		}
	}
	for y := 0; y < res.YLen(); y++ {
		for z := 0; z < res.ZLen(); z++ {
			for x := 0; x < res.XLen(); x++ {
				if res.GetV(x, y, z) == math.MaxUint32 {
					panic("unreachable")
				}
				if res.GetV(x, y, z) > uint32(n2) {
					res.Set(x, y, z, 0)
				}
			}
		}
	}
	return res
}

func WriteNptl(vol BoolVoxelVolume, output io.Writer) (err os.Error) {
	v := 0
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
				if _, err = fmt.Fprintf(output, "%f %f %f %f %f %f\n", float64(x), float64(y), float64(z), nx, ny, nz); err != nil {
					return
				}
			}
		}
	}
	fmt.Fprintf(os.Stderr, "Volume is filled by %v%\n", float64(v)*float64(100)/(float64(vol.XLen())*float64(vol.YLen())*float64(vol.ZLen())))
	return
}

type STLPoint [3]float32

type STLTriangle struct {
	n STLPoint
	v [3]STLPoint
}

func readSTLPoint(a []byte, p *STLPoint) []byte {
	for i := 0; i < 3; i++ {
		cur := uint32(a[0]) + uint32(a[1])<<8 + uint32(a[2])<<16 + uint32(a[3])<<24
		p[i] = *(*float32)(unsafe.Pointer(&cur))
		a = a[4:]
	}
	return a
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

func ReadSTL(r io.Reader) (t []STLTriangle, err os.Error) {
	var data []byte
	if data, err = ioutil.ReadAll(r); err != nil {
		return
	}
	magic := data[:5]
	if string(magic) == "solid" {
		panic("ReadSTL: ascii format is not implemented")
	}
	// Skip STL header
	data = data[80:]
	n := uint32(data[0]) + uint32(data[1])<<8 + uint32(data[2])<<16 + uint32(data[3])<<24
	data = data[4:]

	fmt.Fprintf(os.Stderr, "%d triangles\n", n)
	if len(data) < int(SizeOfSTLTriangle*n) {
		return nil, fmt.Errorf("ReadSTL: unexpected end of file: want %d bytes to read triangle data, but only %d bytes is available", SizeOfSTLTriangle*n, len(data))
	}
	for i := 0; i < int(n); i++ {
		var cur STLTriangle
		data = readSTLPoint(data, &cur.n)
		for j := 0; j < 3; j++ {
			data = readSTLPoint(data, &cur.v[j])
		}
		data = data[2:]
		t = append(t, cur)
	}
	return
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

func Rasterize(m Mesh) (vol BoolVoxelVolume, err os.Error) {
	panic("Rasterize not implemented")
}

func main() {
	triangles, err := ReadSTL(os.Stdin)
	if err != nil {
		log.Fatalf("ReadSTL: %v", err)
	}
	mesh := STLToMesh(VoxelSide*MeshMultiplier, triangles)
	var vol BoolVoxelVolume
	if vol, err = Rasterize(mesh); err != nil {
		log.Fatalf("Rasterize: %v", err)
	}

	/*	vol, err := ReadSchematic(os.Stdin)
		if err != nil {
			log.Fatalf("ReadSchematic: %v", err)
		}*/
	vol = Optimize(vol, 70)
	if err = WriteNptl(vol, os.Stdout); err != nil {
		log.Fatalf("WriteNptl: %v", err)
	}
}
