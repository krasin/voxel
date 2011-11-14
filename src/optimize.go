package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
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
	if _, err = r.r.Read(data); err != nil {
		return
	}
	return string(data), nil
}

func (r *NBTReader) ReadShort() (val int, err os.Error) {
	buf := [2]byte{}
	if _, err = r.r.Read(buf[:]); err != nil {
		return
	}
	val = int(buf[1]) + (int(buf[0]) << 8) // Big Endian
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
	fmt.Printf("Typ: %d, Name: %s\n", typ, name)
	return
}

type SchematicReader struct {
	r *NBTReader
}

func NewSchematicReader(r io.Reader) (sr *SchematicReader, err os.Error) {
	var nr *NBTReader
	if nr, err = NewNBTReader(r); err != nil {
		return
	}
	return &SchematicReader{r: nr}, nil
}

type Schematic struct {
	Width     int
	Length    int
	Height    int
	Materials string
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

func ReadSchematic(input io.Reader) (err os.Error) {
	var r *SchematicReader
	if r, err = NewSchematicReader(input); err != nil {
		return
	}
	_, err = r.Parse()
	return
}

func main() {
	if err := ReadSchematic(os.Stdin); err != nil {
		log.Fatalf("ReadSchematic: %v", err)
	}
}
