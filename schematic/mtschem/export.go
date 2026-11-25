package mtschem

import (
	"blockexchange/parser"
	"blockexchange/types"
	"fmt"
	"io"
)

// Exporter collects schemaparts and writes them as a single mts schematic on Close()
//
// The name-id table precedes the node data in the file, so the whole schematic
// is buffered in memory until Close() is called.
// All nodes are exported with "always place" probability (including air),
// param1 (light) is not exported.
type Exporter struct {
	w         io.Writer
	schematic *Schematic
	nameIDs   map[string]uint16
}

func NewExporter(w io.Writer, sizeX, sizeY, sizeZ int) (*Exporter, error) {
	if sizeX < 1 || sizeY < 1 || sizeZ < 1 || sizeX > MaxSize || sizeY > MaxSize || sizeZ > MaxSize {
		return nil, fmt.Errorf("schematic size %d/%d/%d out of range [1, %d]", sizeX, sizeY, sizeZ, MaxSize)
	}

	return &Exporter{
		w:         w,
		schematic: NewSchematic(sizeX, sizeY, sizeZ),
		nameIDs:   map[string]uint16{"air": 0},
	}, nil
}

func (e *Exporter) Export(schemapart *types.SchemaPart) error {
	mapblock, err := parser.ParseSchemaPart(schemapart)
	if err != nil {
		return err
	}

	// part-local node-id to content-id
	contentIDs := make(map[int16]uint16, len(mapblock.Meta.NodeMapping))
	for name, nodeID := range mapblock.Meta.NodeMapping {
		contentID, found := e.nameIDs[name]
		if !found {
			if len(e.schematic.NodeNames) >= 0xFFFF {
				return fmt.Errorf("too many distinct node names")
			}
			contentID = uint16(len(e.schematic.NodeNames))
			e.schematic.NodeNames = append(e.schematic.NodeNames, name)
			e.nameIDs[name] = contentID
		}
		contentIDs[nodeID] = contentID
	}

	s := e.schematic
	size := mapblock.Meta.Size
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			for z := 0; z < size.Z; z++ {
				sx, sy, sz := schemapart.OffsetX+x, schemapart.OffsetY+y, schemapart.OffsetZ+z
				if sx < 0 || sy < 0 || sz < 0 || sx >= s.SizeX || sy >= s.SizeY || sz >= s.SizeZ {
					return fmt.Errorf("node at %d/%d/%d is outside of the schematic size %d/%d/%d", sx, sy, sz, s.SizeX, s.SizeY, s.SizeZ)
				}

				partIndex := mapblock.GetIndex(x, y, z)
				contentID, found := contentIDs[mapblock.NodeIDS[partIndex]]
				if !found {
					return fmt.Errorf("node-id %d not found in node-mapping", mapblock.NodeIDS[partIndex])
				}

				index := s.Index(sx, sy, sz)
				s.Content[index] = contentID
				s.Param2[index] = mapblock.Param2[partIndex]
			}
		}
	}

	return nil
}

// Close writes the collected schematic to the underlying writer
func (e *Exporter) Close() error {
	return e.schematic.Write(e.w)
}
