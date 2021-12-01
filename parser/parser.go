package parser

import (
	"blockexchange/types"
	"bytes"
	"compress/zlib"
	"encoding/json"
)

const offset = 32768

type SchemaPartSize struct {
	X int
	Y int
	Z int
}

type SchemaPartMetadata struct {
	NodeMapping map[string]int
	Size        SchemaPartSize
}

func getInt(o interface{}) int {
	v, _ := o.(float64)
	return int(v)
}

func (s *SchemaPartMetadata) UnmarshalJSON(data []byte) error {
	m := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	nm := make(map[string]float64)
	err = json.Unmarshal(m["node_mapping"], &nm)
	if err != nil {
		return err
	}

	s.NodeMapping = make(map[string]int)
	for k, v := range nm {
		s.NodeMapping[k] = int(v)
	}

	s.Size = SchemaPartSize{}
	nm = make(map[string]float64)
	err = json.Unmarshal(m["size"], &nm)
	if err != nil {
		return err
	}

	s.Size.X = getInt(nm["x"])
	s.Size.Y = getInt(nm["y"])
	s.Size.Z = getInt(nm["z"])

	return nil
}

type ParsedSchemaPart struct {
	NodeIDS []int16
	Param1  []byte
	Param2  []byte
	Meta    *SchemaPartMetadata
}

func (mapblock *ParsedSchemaPart) GetIndex(x, y, z int) int {
	return z + (y * mapblock.Meta.Size.Z) + (x * mapblock.Meta.Size.Y * mapblock.Meta.Size.Z)
}

func ParseSchemaPart(part *types.SchemaPart) (*ParsedSchemaPart, error) {
	r, err := zlib.NewReader(bytes.NewReader(part.MetaData))
	if err != nil {
		return nil, err
	}

	md := SchemaPartMetadata{}
	err = json.NewDecoder(r).Decode(&md)
	if err != nil {
		return nil, err
	}

	r, err = zlib.NewReader(bytes.NewReader(part.Data))
	if err != nil {
		return nil, err
	}

	result := ParsedSchemaPart{
		Meta: &md,
	}

	size := md.Size.X * md.Size.Y * md.Size.Z

	// prepare result buffer
	result.NodeIDS = make([]int16, size)
	result.Param1 = make([]byte, size)
	result.Param2 = make([]byte, size)

	data := make([]byte, size*4)
	r.Read(data)

	for i := 0; i < size; i++ {
		result.NodeIDS[i] = int16((int(data[i*2]) * 256) + int(data[(i*2)+1]) - offset)
		result.Param1[i] = data[(size*2)+i]
		result.Param2[i] = data[(size*3)+i]
	}

	return &result, nil
}
