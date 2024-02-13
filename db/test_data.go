package db

import (
	"blockexchange/types"
	"time"

	"github.com/vingarcia/ksql"
	"golang.org/x/crypto/bcrypt"
)

func PopulateTestData(kdb ksql.Provider) error {
	userrepo := &UserRepository{kdb: kdb}
	tokenrepo := &AccessTokenRepository{kdb: kdb}

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
		UserUID:  user.UID,
		Created:  time.Now().Unix() * 1000,
		Expires:  (time.Now().Unix() * 1000) + (1000 * 3600 * 24 * 30 * 365),
		UseCount: 0,
	}

	err = tokenrepo.CreateAccessToken(token)
	return err
}
