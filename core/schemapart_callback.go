package core

import "blockexchange/types"

func (c *Core) SchemapartCallback(schema_id int64, cb func(sp *types.SchemaPart) error) error {

	schemapart, err := c.repos.SchemaPartRepo.GetFirstBySchemaID(schema_id)
	if err != nil {
		return err
	}
	err = cb(schemapart)
	if err != nil {
		return err
	}

	for {
		schemapart, err = c.repos.SchemaPartRepo.GetNextBySchemaIDAndOffset(schema_id, schemapart.OffsetX, schemapart.OffsetY, schemapart.OffsetZ)
		if err != nil {
			return err
		}

		if schemapart != nil {
			err = cb(schemapart)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}
