package parser

type SchemaPartSize struct {
	X int
	Y int
	Z int
}

type Inventory struct {
}

type Fields struct {
}

type MetadataEntry struct {
	Inventories map[string]*Inventory
	Fields      map[string]*Fields
}

type Metadata struct {
	Meta *MetadataEntry
	//timers
}

type SchemaPartMetadata struct {
	NodeMapping map[string]int
	Size        SchemaPartSize
	Metadata    *Metadata
}

type ParsedSchemaPart struct {
	NodeIDS []int16
	Param1  []byte
	Param2  []byte
	Meta    *SchemaPartMetadata
}
