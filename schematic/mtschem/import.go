package mtschem

import (
	"blockexchange/parser"
)

// Import splits the schematic into mapblock-sized schemaparts
//
// Lossy conversions:
//   - nodes with a placement probability of 0 ("never place", usually air) are imported as "air"
//   - param1 (probability and force-place bit) is not imported and set to 0
//   - y-slice probabilities are ignored
//
// Parts without any non-air node are omitted.
func Import(s *Schematic) ([]*parser.ParsedSchemaPart, error) {
	parts := make([]*parser.ParsedSchemaPart, 0)

	for mbx := 0; mbx*16 < s.SizeX; mbx++ {
		for mby := 0; mby*16 < s.SizeY; mby++ {
			for mbz := 0; mbz*16 < s.SizeZ; mbz++ {
				part := importPart(s, mbx, mby, mbz)
				if part != nil {
					parts = append(parts, part)
				}
			}
		}
	}

	return parts, nil
}

func importPart(s *Schematic, mbx, mby, mbz int) *parser.ParsedSchemaPart {
	size := parser.SchemaPartSize{
		X: min(16, s.SizeX-mbx*16),
		Y: min(16, s.SizeY-mby*16),
		Z: min(16, s.SizeZ-mbz*16),
	}
	bufsize := size.X * size.Y * size.Z

	part := &parser.ParsedSchemaPart{
		Meta: &parser.SchemaPartMetadata{
			NodeMapping: map[string]int16{"air": 0},
			Size:        size,
			Metadata: &parser.Metadata{
				Meta: make(map[string]*parser.MetadataEntry),
			},
		},
		NodeIDS: make([]int16, bufsize),
		Param1:  make([]byte, bufsize),
		Param2:  make([]byte, bufsize),
		PosX:    mbx,
		PosY:    mby,
		PosZ:    mbz,
	}

	// content-id to part-local node-id (max 16^3 nodes + air per part, fits into int16)
	nodeIDs := map[uint16]int16{}
	hasNodes := false

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			for z := 0; z < size.Z; z++ {
				index := s.Index(mbx*16+x, mby*16+y, mbz*16+z)
				if s.Param1[index]&ProbMask == ProbNever {
					// never placed, keep air
					continue
				}

				contentID := s.Content[index]
				name := s.NodeNames[contentID]
				if name == "air" {
					continue
				}

				nodeID, found := nodeIDs[contentID]
				if !found {
					nodeID, found = part.Meta.NodeMapping[name]
					if !found {
						nodeID = int16(len(part.Meta.NodeMapping))
						part.Meta.NodeMapping[name] = nodeID
					}
					nodeIDs[contentID] = nodeID
				}

				partIndex := part.GetIndex(x, y, z)
				part.NodeIDS[partIndex] = nodeID
				part.Param2[partIndex] = s.Param2[index]
				hasNodes = true
			}
		}
	}

	if !hasNodes {
		return nil
	}
	return part
}
