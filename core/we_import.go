package core

import (
	"blockexchange/schematic/worldedit"
	"blockexchange/types"
	"fmt"
	"time"
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

	newSchemaName, err := c.FindUnusedSchemaname(schemaname, username)
	if err != nil {
		return nil, err
	}

	user, err := c.repos.UserRepo.GetUserByName(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: '%s'", username)
	}

	schema := &types.Schema{
		Created:     time.Now().Unix() * 1000,
		UserID:      *user.ID,
		Name:        newSchemaName,
		Mtime:       time.Now().Unix() * 1000,
		Description: "WE Import",
		Complete:    false,
		SizeX:       max_x + 1,
		SizeY:       max_y + 1,
		SizeZ:       max_z + 1,
		TotalParts:  len(parts),
		License:     "CC0",
	}

	err = c.repos.SchemaRepo.CreateSchema(schema)
	if err != nil {
		return nil, err
	}

	for _, part := range parts {
		sp, err := part.Convert()
		if err != nil {
			return nil, err
		}
		sp.SchemaID = *schema.ID

		err = c.repos.SchemaPartRepo.CreateOrUpdateSchemaPart(sp)
		if err != nil {
			return nil, err
		}
	}

	for _, modname := range modnames {
		err = c.repos.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
			SchemaID: *schema.ID,
			ModName:  modname,
		})
		if err != nil {
			return nil, err
		}
	}

	return c.repos.SchemaRepo.GetSchemaById(*schema.ID)
}
