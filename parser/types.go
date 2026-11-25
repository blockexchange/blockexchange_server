package parser

import (
	"blockexchange/types"
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"
)

// zlibWriterPool reuses zlib writers, a new writer allocates around 1 MB of compressor state
var zlibWriterPool = sync.Pool{
	New: func() any { return zlib.NewWriter(nil) },
}

// compress returns the zlib-compressed data
func compress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, len(data)/2+64))
	w := zlibWriterPool.Get().(*zlib.Writer)
	defer zlibWriterPool.Put(w)
	w.Reset(buf)

	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type SchemaPartSize struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type MetadataEntry struct {
	Inventories map[string][]string `json:"inventory"`
	Fields      map[string]string   `json:"fields"`
}

type MetadataTimer struct {
	Timeout float64 `json:"timeout"`
	Elapsed float64 `json:"elapsed"`
}

type Metadata struct {
	Meta   map[string]*MetadataEntry `json:"meta"`
	Timers map[string]*MetadataTimer `json:"timers"`
}

func (m Metadata) GetKey(x, y, z int) string {
	return fmt.Sprintf("(%d,%d,%d)", x, y, z)
}

type SchemaPartMetadata struct {
	NodeMapping map[string]int16 `json:"node_mapping"`
	Size        SchemaPartSize   `json:"size"`
	Metadata    *Metadata        `json:"metadata"`
}

type ParsedSchemaPart struct {
	NodeIDS        []int16
	Param1         []byte
	Param2         []byte
	PosX           int
	PosY           int
	PosZ           int
	Meta           *SchemaPartMetadata
	NodeNameLookup map[int16]string
}

func (psp *ParsedSchemaPart) GetIndex(x, y, z int) int {
	return z + (y * psp.Meta.Size.Z) + (x * psp.Meta.Size.Y * psp.Meta.Size.Z)
}

// IsAirOnly returns true if the part consists only of air nodes
func (psp *ParsedSchemaPart) IsAirOnly() bool {
	_, found := psp.Meta.NodeMapping["air"]
	return found && len(psp.Meta.NodeMapping) == 1
}

func (psp *ParsedSchemaPart) Convert() (*types.SchemaPart, error) {
	sp := &types.SchemaPart{
		OffsetX: psp.PosX * 16,
		OffsetY: psp.PosY * 16,
		OffsetZ: psp.PosZ * 16,
		Mtime:   time.Now().UnixMilli(),
	}
	rawMeta, err := json.Marshal(psp.Meta)
	if err != nil {
		return nil, err
	}

	sp.MetaData, err = compress(rawMeta)
	if err != nil {
		return nil, err
	}

	size := psp.Meta.Size.X * psp.Meta.Size.Y * psp.Meta.Size.Z
	mapdata := make([]byte, size*4)

	for i := 0; i < size; i++ {
		//TODO: document this monstrosity
		mapdata[i*2] = byte(int(math.Floor((float64(psp.NodeIDS[i])+32768)/256)) % 256)
		mapdata[(i*2)+1] = byte((int(float64(psp.NodeIDS[i]) + 32768)) % 256)
		mapdata[(size*2)+i] = psp.Param1[i]
		mapdata[(size*3)+i] = psp.Param2[i]
	}

	sp.Data, err = compress(mapdata)
	if err != nil {
		return nil, err
	}

	return sp, nil
}
