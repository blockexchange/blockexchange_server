package db

import (
	"blockexchange/types"
)

func GetSchemaById(id int64) (*types.Schema, error) {
	schema := types.Schema{}
	query := `
		select
			id, created, user_id, name,
			max_x, max_y, max_z
		from schema where id = $1
	`
	row := DB.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&schema.ID, &schema.Created, &schema.UserID, &schema.Name,
		&schema.MaxX, &schema.MaxY, &schema.MaxZ,
	)

	return &schema, err
}
