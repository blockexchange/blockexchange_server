package core

import (
	"blockexchange/schematic/bx"
	"blockexchange/types"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (c *Core) ImportBX(data []byte, username string) (*types.Schema, error) {
	res, err := bx.ParseBXContent(data)
	if err != nil {
		return nil, err
	}

	newSchemaName, err := c.FindUnusedSchemaname(res.Schema.Name, username)
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

	res.Schema.UID = uuid.NewString()
	res.Schema.Created = time.Now().UnixMilli()
	res.Schema.Mtime = time.Now().UnixMilli()
	res.Schema.UserUID = user.UID
	res.Schema.Name = newSchemaName

	err = c.repos.SchemaRepo.CreateSchema(res.Schema)
	if err != nil {
		return nil, err
	}

	for _, modname := range res.Mods {
		err = c.repos.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
			SchemaUID: res.Schema.UID,
			ModName:   modname,
		})
		if err != nil {
			return nil, err
		}
	}

	for _, part := range res.Parts {
		part.SchemaUID = res.Schema.UID
		part.Mtime = res.Schema.Mtime
		err = c.repos.SchemaPartRepo.CreateOrUpdateSchemaPart(part)
		if err != nil {
			return nil, err
		}
	}

	return res.Schema, nil
}
