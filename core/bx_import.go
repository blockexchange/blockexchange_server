package core

import (
	"archive/zip"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"time"
)

type BXImport struct {
	Parts  []*types.SchemaPart
	Schema *types.Schema
	Mods   []string
}

func ParseBXContent(data []byte) (*BXImport, error) {
	res := &BXImport{
		Parts:  make([]*types.SchemaPart, 0),
		Schema: &types.Schema{},
		Mods:   make([]string, 0),
	}

	r := bytes.NewReader(data)

	zr, err := zip.NewReader(r, int64(len(data)))
	if err != nil {
		return nil, err
	}

	err = unmarshalFromZip(zr, "schema.json", res.Schema)
	if err != nil {
		return nil, err
	}

	err = unmarshalFromZip(zr, "mods.json", &res.Mods)
	if err != nil {
		return nil, err
	}

	mb_max_x := int(res.Schema.SizeX / 16)
	mb_max_y := int(res.Schema.SizeY / 16)
	mb_max_z := int(res.Schema.SizeZ / 16)
	for x := 0; x <= mb_max_x; x++ {
		for y := 0; y <= mb_max_y; y++ {
			for z := 0; z <= mb_max_z; z++ {
				sp := &types.SchemaPart{
					OffsetX: x * 16,
					OffsetY: y * 16,
					OffsetZ: z * 16,
				}
				err = unmarshalFromZip(zr, formatSchemapartFilename(sp), sp)
				if errors.Is(err, fs.ErrNotExist) {
					// air-only mapblock
					continue
				}
				if err != nil {
					// other error
					return nil, err
				}

				res.Parts = append(res.Parts, sp)
			}
		}
	}

	return res, nil
}

func unmarshalFromZip(zr *zip.Reader, name string, obj any) error {
	f, err := zr.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(obj)
}

func (c *Core) ImportBX(data []byte, username string) (*types.Schema, error) {
	res, err := ParseBXContent(data)
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

	res.Schema.Created = time.Now().UnixMilli()
	res.Schema.Mtime = time.Now().UnixMilli()
	res.Schema.UserID = *user.ID
	res.Schema.Name = newSchemaName

	err = c.repos.SchemaRepo.CreateSchema(res.Schema)
	if err != nil {
		return nil, err
	}

	for _, modname := range res.Mods {
		err = c.repos.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
			SchemaID: res.Schema.ID,
			ModName:  modname,
		})
		if err != nil {
			return nil, err
		}
	}

	for _, part := range res.Parts {
		part.SchemaID = res.Schema.ID
		part.Mtime = res.Schema.Mtime
		err = c.repos.SchemaPartRepo.CreateOrUpdateSchemaPart(part)
		if err != nil {
			return nil, err
		}
	}

	return res.Schema, nil
}
