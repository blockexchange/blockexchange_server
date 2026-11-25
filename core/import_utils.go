package core

import (
	"blockexchange/parser"
	"blockexchange/types"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// importedSchema holds the result of a parsed schematic file (worldedit, mts)
type importedSchema struct {
	Username    string
	Schemaname  string
	Description string
	SizeX       int
	SizeY       int
	SizeZ       int
	Parts       []*parser.ParsedSchemaPart
	Modnames    []string
}

// modnamesFromParts returns the sorted mod names of all node names in the parts
func modnamesFromParts(parts []*parser.ParsedSchemaPart) []string {
	modnames := map[string]bool{}
	for _, part := range parts {
		for name := range part.Meta.NodeMapping {
			modname, _, found := strings.Cut(name, ":")
			if found && modname != "" {
				modnames[modname] = true
			}
		}
	}
	res := make([]string, 0, len(modnames))
	for modname := range modnames {
		res = append(res, modname)
	}
	sort.Strings(res)
	return res
}

// createImportedSchema stores a new (incomplete) schema with its parts and mods, PostImport() completes it
func (c *Core) createImportedSchema(is *importedSchema) (*types.Schema, error) {
	newSchemaName, err := c.FindUnusedSchemaname(is.Schemaname, is.Username)
	if err != nil {
		return nil, err
	}

	user, err := c.repos.UserRepo.GetUserByName(is.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found: '%s'", is.Username)
	}

	schema := &types.Schema{
		UID:         uuid.NewString(),
		Created:     time.Now().UnixMilli(),
		UserUID:     user.UID,
		Name:        newSchemaName,
		Mtime:       time.Now().UnixMilli(),
		Description: is.Description,
		Complete:    false,
		SizeX:       is.SizeX,
		SizeY:       is.SizeY,
		SizeZ:       is.SizeZ,
		TotalParts:  len(is.Parts),
		License:     "CC0",
	}

	err = c.repos.SchemaRepo.CreateSchema(schema)
	if err != nil {
		return nil, err
	}

	for _, part := range is.Parts {
		sp, err := part.Convert()
		if err != nil {
			return nil, err
		}
		sp.SchemaUID = schema.UID

		err = c.repos.SchemaPartRepo.CreateOrUpdateSchemaPart(sp)
		if err != nil {
			return nil, err
		}
	}

	for _, modname := range is.Modnames {
		err = c.repos.SchemaModRepo.CreateSchemaMod(&types.SchemaMod{
			SchemaUID: schema.UID,
			ModName:   modname,
		})
		if err != nil {
			return nil, err
		}
	}

	return c.repos.SchemaRepo.GetSchemaByUID(schema.UID)
}

func (c *Core) PostImport(schema *types.Schema) (*types.Schema, error) {
	schema.Complete = true
	err := c.repos.SchemaRepo.UpdateSchema(schema)
	if err != nil {
		return nil, fmt.Errorf("update error: %v", err)
	}

	err = c.repos.SchemaRepo.CalculateStats(schema.UID)
	if err != nil {
		return nil, fmt.Errorf("stats calc error: %v", err)
	}

	_, err = c.UpdatePreview(schema)
	if err != nil {
		return nil, fmt.Errorf("preview update error: %v", err)
	}

	return c.repos.SchemaRepo.GetSchemaByUID(schema.UID)
}

func (c *Core) FindUnusedSchemaname(schemaname, username string) (string, error) {
	if schemaname == "" || !ValidateName(schemaname) {
		// placeholder
		return fmt.Sprintf("import_%d", rand.Int()), nil
	}

	newSchemaName := schemaname
	existing_schema, err := c.repos.SchemaRepo.GetSchemaByUsernameAndName(username, newSchemaName)
	if err != nil {
		return "", fmt.Errorf("GetSchemaByUsernameAndName error: %v", err)
	}
	if existing_schema == nil {
		// no previous schema and valid name
		return newSchemaName, nil
	}

	// schema with that name already exists, add number to name
	i := 2
	for {
		newSchemaName = fmt.Sprintf("%s_%d", schemaname, i)
		existing_schema, err = c.repos.SchemaRepo.GetSchemaByUsernameAndName(username, newSchemaName)
		if err != nil {
			return "", fmt.Errorf("GetSchemaByUsernameAndName error: %v", err)
		}
		if existing_schema == nil {
			return newSchemaName, nil
		}
		i++
		if i > 50 {
			break
		}
	}

	// nothing helped, fall back to random name
	return fmt.Sprintf("import_%d", rand.Int()), nil
}
