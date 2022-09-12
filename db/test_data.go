package db

import (
	"blockexchange/types"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func PopulateTestData(_db *sqlx.DB) error {
	userrepo := DBUserRepository{DB: _db}
	tokenrepo := AccessTokenRepository{DB: _db.DB}

	user, err := userrepo.GetUserByName("Testuser")
	if err != nil {
		return err
	}
	if user != nil {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}

	user = &types.User{
		Name: "Testuser",
		Type: types.UserTypeLocal,
		Hash: string(hash),
	}
	err = userrepo.CreateUser(user)
	if err != nil {
		return err
	}

	token := &types.AccessToken{
		Name:     "Default",
		Token:    "default",
		UserID:   user.ID,
		Created:  time.Now().Unix() * 1000,
		Expires:  (time.Now().Unix() * 1000) + (1000 * 3600 * 24 * 30 * 365),
		UseCount: 0,
	}

	err = tokenrepo.CreateAccessToken(token)
	return err
}
