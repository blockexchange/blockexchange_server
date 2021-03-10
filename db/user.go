package db

import (
	"blockexchange/types"
	"database/sql"
)

var fields = "id, created, name, hash, type, external_id, mail"

func mapFromDB(row *sql.Row) (*types.User, error) {
	user := types.User{}

	err := row.Scan(
		&user.ID, &user.Created, &user.Name,
		&user.Hash, &user.Type, &user.ExternalID, &user.Mail,
	)

	return &user, err
}

func GetUserById(id int64) (*types.User, error) {
	query := `
		select ` + fields + `
		from public.user where id = $1
	`
	row := DB.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return mapFromDB(row)
}
