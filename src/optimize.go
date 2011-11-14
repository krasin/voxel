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
	buf := [2]byte{}
	if _, err = r.r.Read(buf[:]); err != nil {
		return
	}
	l := int(buf[1]) + (int(buf[0]) << 8) // Big Endian
	data := make([]byte, l)
	if _, err = r.r.Read(data); err != nil {
		return
	}
	return string(data), nil
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

func main() {
	r, err := NewNBTReader(os.Stdin)
	if err != nil {
		log.Fatalf("gzip: %v", err)
	}
	var typ byte
	var name string
	if typ, name, err = r.ReadTagName(); err != nil {
		log.Fatalf("ReadTagName: %v", err)
	}
	fmt.Printf("Typ: %d, Name: %s\n", typ, name)
}
