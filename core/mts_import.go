package core

import (
	"blockexchange/schematic/mtschem"
	"blockexchange/types"
)

func (c *Core) ImportMTS(data []byte, username, schemaname string) (*types.Schema, error) {
	s, err := mtschem.Parse(data)
	if err != nil {
		return nil, err
	}

	parts, err := mtschem.Import(s)
	if err != nil {
		return nil, err
	}

	return c.createImportedSchema(&importedSchema{
		Username:    username,
		Schemaname:  schemaname,
		Description: "MTS Import",
		SizeX:       s.SizeX,
		SizeY:       s.SizeY,
		SizeZ:       s.SizeZ,
		Parts:       parts,
		Modnames:    modnamesFromParts(parts),
	})
}
