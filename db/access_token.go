package db

import "blockexchange/types"

func GetAccessTokensByUserID(user_id int64) ([]types.AccessToken, error) {
	list := []types.AccessToken{}
	err := DB.Select(&list, "select * from access_token where user_id = $1", user_id)
	if err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func CreateAccessToken(access_token *types.AccessToken) error {
	query := `
		insert into
		access_token(
			name, token, user_id,
			created, expires
		)
		values(
			$1, $2, $3,
			$4, $5
		)
		returning id
	`
	row := DB.QueryRow(query,
		access_token.Name, access_token.Token, access_token.UserID,
		access_token.Created, access_token.Expires,
	)
	return row.Scan(&access_token.ID)
}

func IncrementAccessTokenUseCount(id int64) error {
	_, err := DB.Exec("update access_token set usecount = usecount + 1 where id = $1", id)
	return err
}

func RemoveAccessToken(id, user_id int64) error {
	_, err := DB.Exec("delete from access_token where id = $1 and user_id = $2", id, user_id)
	return err
}
