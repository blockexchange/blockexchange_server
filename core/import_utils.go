package core

import (
	"fmt"
	"math/rand"
)

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
	if existing_schema != nil {
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
	}

	return fmt.Sprintf("import_%d", rand.Int()), nil
}
