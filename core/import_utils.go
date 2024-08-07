package core

import (
	"blockexchange/types"
	"fmt"
	"math/rand"
)

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
