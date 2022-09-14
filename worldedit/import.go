package worldedit

import (
	"blockexchange/parser"
	"fmt"
	"math"
)

func getMapblockPos(x, y, z int) (int, int, int) {
	return int(math.Floor(float64(x) / 16)), int(math.Floor(float64(y) / 16)), int(math.Floor(float64(z) / 16))
}

func formatPos(x, y, z int) string {
	return fmt.Sprintf("%d/%d/%d", x, y, z)
}

func Import(entries []*WEEntry) ([]*parser.ParsedSchemaPart, error) {
	max_x, max_y, max_z := GetBoundaries(entries)
	max_mbx, max_mby, max_mbz := getMapblockPos(max_x, max_y, max_z)
	parts := make(map[string]*parser.ParsedSchemaPart)
	var nextNodeID int16 = 1

	for _, entry := range entries {
		mbx, mby, mbz := getMapblockPos(entry.X, entry.Y, entry.Z)
		index := formatPos(mbx, mby, mbz)

		// intra-mapblock coords
		rel_x := entry.X - (mbx * 16)
		rel_y := entry.Y - (mby * 16)
		rel_z := entry.Z - (mbz * 16)

		part := parts[index]
		if part == nil {
			// create a new part
			part = &parser.ParsedSchemaPart{
				Meta: &parser.SchemaPartMetadata{
					NodeMapping: make(map[string]int16),
					Size:        parser.SchemaPartSize{},
					Metadata: &parser.Metadata{
						Meta: &parser.MetadataEntry{
							Fields:      make(map[string]parser.Fields),
							Inventories: make(map[string]parser.Inventory),
						},
					},
				},
				PosX: mbx,
				PosY: mby,
				PosZ: mbz,
			}

			// calculate size
			if mbx < max_mbx {
				part.Meta.Size.X = 16
			} else {
				part.Meta.Size.X = (max_x % 16) + 1
			}
			if mby < max_mby {
				part.Meta.Size.Y = 16
			} else {
				part.Meta.Size.Y = (max_y % 16) + 1
			}
			if mbz < max_mbz {
				part.Meta.Size.Z = 16
			} else {
				part.Meta.Size.Z = (max_z % 16) + 1
			}

			// populate mapdata
			bufsize := part.Meta.Size.X * part.Meta.Size.Y * part.Meta.Size.Z
			part.NodeIDS = make([]int16, bufsize)
			part.Param1 = make([]byte, bufsize)
			part.Param2 = make([]byte, bufsize)

			// set air nodeid
			part.Meta.NodeMapping["air"] = 0
		}
		parts[index] = part

		var nodeID int16 = part.Meta.NodeMapping[entry.Name]
		if nodeID == 0 {
			nodeID = nextNodeID
			nextNodeID++
			part.Meta.NodeMapping[entry.Name] = nodeID
		}

		bufIndex := part.GetIndex(rel_x, rel_y, rel_z)
		part.NodeIDS[bufIndex] = nodeID
		part.Param1[bufIndex] = entry.Param1
		part.Param2[bufIndex] = entry.Param2

		if entry.Meta != nil {
			key := part.Meta.Metadata.Meta.GetKey(rel_x, rel_y, rel_z)
			part.Meta.Metadata.Meta.Fields[key] = parser.Fields(entry.Meta.Fields)
			part.Meta.Metadata.Meta.Inventories[key] = parser.Inventory(entry.Meta.Inventory)
		}
	}

	res := make([]*parser.ParsedSchemaPart, 0)
	for _, part := range parts {
		res = append(res, part)
	}

	return res, nil
}
