package db

import (
	"blockexchange/types"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func PopulateTestData(repos *Repositories) error {

	user, err := repos.UserRepo.GetUserByName("Testuser")
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
	err = repos.UserRepo.CreateUser(user)
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

	err = repos.AccessTokenRepo.CreateAccessToken(token)
	return err
}
