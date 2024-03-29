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
		return nil, fmt.Errorf("parse error: %v", err)
	}

	newSchemaName, err := c.FindUnusedSchemaname(res.Schema.Name, username)
	if err != nil {
		return nil, fmt.Errorf("find unused schemaname error: %v", err)
	}

	user, err := c.repos.UserRepo.GetUserByName(username)
	if err != nil {
		return nil, fmt.Errorf("get user by name error: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: '%s'", username)
	}

	res.Schema.UID = uuid.NewString()
	res.Schema.Downloads = 0
	res.Schema.Views = 0
	res.Schema.Created = time.Now().UnixMilli()
	res.Schema.Mtime = time.Now().UnixMilli()
	res.Schema.UserUID = user.UID
	res.Schema.Name = newSchemaName

	err = c.repos.SchemaRepo.CreateSchema(res.Schema)
	if err != nil {
		return nil, fmt.Errorf("create schema error: %v", err)
	}

	for _, modname := range res.Mods {
		err = c.repos.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
			SchemaUID: res.Schema.UID,
			ModName:   modname,
		})
		if err != nil {
			return nil, fmt.Errorf("create schemamod error: %v", err)
		}
	}

	for _, part := range res.Parts {
		part.SchemaUID = res.Schema.UID
		part.Mtime = res.Schema.Mtime
		err = c.repos.SchemaPartRepo.CreateOrUpdateSchemaPart(part)
		if err != nil {
			return nil, fmt.Errorf("create or update schemapart error: %v", err)
		}
	}

	return res.Schema, nil
}
