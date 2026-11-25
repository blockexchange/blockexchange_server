package core

import (
	"blockexchange/schematic/worldedit"
	"blockexchange/types"
)

func (c *Core) ImportWE(data []byte, username, schemaname string) (*types.Schema, error) {
	entries, modnames, err := worldedit.Parse(data)
	if err != nil {
		return nil, err
	}
	max_x, max_y, max_z := worldedit.GetBoundaries(entries)

	parts, err := worldedit.Import(entries)
	if err != nil {
		return nil, err
	}

	return c.createImportedSchema(&importedSchema{
		Username:    username,
		Schemaname:  schemaname,
		Description: "WE Import",
		SizeX:       max_x + 1,
		SizeY:       max_y + 1,
		SizeZ:       max_z + 1,
		Parts:       parts,
		Modnames:    modnames,
	})
}
