package testdata

import (
	"blockexchange/types"
	"errors"
)

type MockUserRepository struct {
	Users []types.User
}

func (repo MockUserRepository) GetUserById(id int64) (*types.User, error) {
	return nil, errors.New("not implemented")
}
func (repo MockUserRepository) GetUserByName(name string) (*types.User, error) {
	for _, user := range repo.Users {
		if user.Name == name {
			return &user, nil
		}
	}
	return nil, nil
}
func (repo MockUserRepository) GetUserByExternalId(external_id string) (*types.User, error) {
	return nil, errors.New("not implemented")
}
func (repo MockUserRepository) GetUsers() ([]types.User, error) {
	return nil, errors.New("not implemented")
}
func (repo MockUserRepository) CreateUser(user *types.User) error {
	return errors.New("not implemented")
}
func (repo MockUserRepository) UpdateUser(user *types.User) error {
	return errors.New("not implemented")
}
