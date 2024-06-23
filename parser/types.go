package parser

import (
	"blockexchange/types"
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type SchemaPartSize struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type MetadataEntry struct {
	Inventories map[string][]string `json:"inventory"`
	Fields      map[string]string   `json:"fields"`
}

type Metadata struct {
	Meta map[string]*MetadataEntry `json:"meta"`
	//TODO: timers
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

func (psp *ParsedSchemaPart) Convert() (*types.SchemaPart, error) {
	sp := &types.SchemaPart{
		OffsetX: psp.PosX * 16,
		OffsetY: psp.PosY * 16,
		OffsetZ: psp.PosZ * 16,
		Mtime:   time.Now().Unix() * 1000,
	}
	rawMeta, err := json.Marshal(psp.Meta)
	if err != nil {
		return nil, err
	}

	metaBuf := bytes.NewBuffer([]byte{})
	w := zlib.NewWriter(metaBuf)
	_, err = w.Write(rawMeta)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	sp.MetaData = metaBuf.Bytes()

	size := psp.Meta.Size.X * psp.Meta.Size.Y * psp.Meta.Size.Z
	mapdata := make([]byte, size*4)

	for i := 0; i < size; i++ {
		//TODO: document this monstrosity
		mapdata[i*2] = byte(int(math.Floor((float64(psp.NodeIDS[i])+32768)/256)) % 256)
		mapdata[(i*2)+1] = byte((int(float64(psp.NodeIDS[i]) + 32768)) % 256)
		mapdata[(size*2)+i] = psp.Param1[i]
		mapdata[(size*3)+i] = psp.Param2[i]
	}

	mapdataBuf := bytes.NewBuffer([]byte{})
	w = zlib.NewWriter(mapdataBuf)
	_, err = w.Write(mapdata)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	sp.Data = mapdataBuf.Bytes()

	return sp, nil
}
