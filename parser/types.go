package parser

import "fmt"

type SchemaPartSize struct {
	X int
	Y int
	Z int
}

type Inventory map[string][]string
type Fields map[string]string

type MetadataEntry struct {
	Inventories map[string]Inventory `json:"inventory"`
	Fields      map[string]Fields    `json:"fields"`
}

func (m MetadataEntry) GetKey(x, y, z int) string {
	return fmt.Sprintf("(%d,%d,%d)", x, y, z)
}

type Metadata struct {
	Meta *MetadataEntry `json:"meta"`
	//TODO: timers
}

type SchemaPartMetadata struct {
	NodeMapping map[string]int16 `json:"node_mapping"`
	Size        SchemaPartSize   `json:"size"`
	Metadata    *Metadata        `json:"metadata"`
}

type ParsedSchemaPart struct {
	NodeIDS []int16
	Param1  []byte
	Param2  []byte
	PosX    int
	PosY    int
	PosZ    int
	Meta    *SchemaPartMetadata
}

func (mapblock *ParsedSchemaPart) GetIndex(x, y, z int) int {
	return z + (y * mapblock.Meta.Size.Z) + (x * mapblock.Meta.Size.Y * mapblock.Meta.Size.Z)
}
