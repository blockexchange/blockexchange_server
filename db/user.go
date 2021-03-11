package db

import (
	"blockexchange/types"
	"database/sql"

	"github.com/sirupsen/logrus"
)

var fields = "id, created, name, hash, type, external_id, mail"

func mapFromDB(row Scanner) (*types.User, error) {
	if row.Err() != nil {
		logrus.Error(row.Err())
		return nil, row.Err()
	}

	user := types.User{}

	err := row.Scan(
		&user.ID, &user.Created, &user.Name,
		&user.Hash, &user.Type, &user.ExternalID, &user.Mail,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func GetUserById(id int64) (*types.User, error) {
	logrus.Trace("db.GetUserById", id)
	query := "select " + fields + " from public.user where id = $1"
	row := DB.QueryRow(query, id)
	return mapFromDB(row)
}

func GetUserByExternalId(external_id string) (*types.User, error) {
	logrus.Trace("db.GetUserByExternalId", external_id)
	query := "select " + fields + " from public.user where external_id = $1"
	row := DB.QueryRow(query, external_id)
	return mapFromDB(row)
}

func GetUsers() ([]*types.User, error) {
	logrus.Trace("db.GetUsers")
	query := "select " + fields + " from public.user"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	users := make([]*types.User, 0)
	for rows.Next() {
		user, err := mapFromDB(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func CreateUser(user *types.User) error {
	logrus.Trace("db.CreateUser", user)
	query := `
		insert into
		public.user(
			created, name, hash, type,
			external_id, mail
		)
		values(
			$1, $2, $3, $4,
			$5, $6
		)
		returning id
	`
	row := DB.QueryRow(query,
		user.Created, user.Name, user.Hash, user.Type,
		user.ExternalID, user.Mail,
	)
	return row.Scan(&user.ID)
}

func UpdateUser(user *types.User) error {
	logrus.Trace("db.UpdateUser", user)
	query := `
	update public.user
	where id = $1
	set name = $2, mail = $3
	`
	_, err := DB.Exec(query, user.ID, user.Name, user.Mail)
	return err
}
