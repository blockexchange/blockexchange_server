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
	//max_x, max_y, max_z := GetBoundaries(entries)
	parts := make(map[string]*parser.ParsedSchemaPart)

	for _, entry := range entries {
		mbx, mby, mbz := getMapblockPos(entry.X, entry.Y, entry.Z)
		index := formatPos(mbx, mby, mbz)
		part := parts[index]
		if part == nil {
			part = &parser.ParsedSchemaPart{
				Meta: &parser.SchemaPartMetadata{
					Size: parser.SchemaPartSize{},
				},
			}
		}
		parts[index] = part
	}

	res := make([]*parser.ParsedSchemaPart, 0)
	for _, part := range parts {
		res = append(res, part)
		// TODO: mapbloc pos
	}

	return res, nil
}
