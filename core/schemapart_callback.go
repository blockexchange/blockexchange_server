package core

import (
	"blockexchange/types"
	"fmt"
)

func (c *Core) SchemapartCallback(schema_uid string, cb func(sp *types.SchemaPart) error) error {

	schemapart, err := c.repos.SchemaPartRepo.GetFirstBySchemaUID(schema_uid)
	if err != nil {
		return fmt.Errorf("get first schema error: %v", err)
	}
	err = cb(schemapart)
	if err != nil {
		return fmt.Errorf("callback error: %v", err)
	}

	for {
		schemapart, err = c.repos.SchemaPartRepo.GetNextBySchemaUIDAndOffset(schema_uid, schemapart.OffsetX, schemapart.OffsetY, schemapart.OffsetZ)
		if err != nil {
			return fmt.Errorf("get next schema error: %v", err)
		}

		if schemapart != nil {
			err = cb(schemapart)
			if err != nil {
				return fmt.Errorf("next callback error: %v", err)
			}
		} else {
			break
		}
	}

	return nil
}
