package db

import (
	"blockexchange/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	g *gorm.DB
}

func (r *UserRepository) GetUserByUID(uid string) (*types.User, error) {
	return FindSingle[types.User](r.g.Where(types.User{UID: uid}))
}

func (r *UserRepository) CountUsers() (int64, error) {
	var c int64
	return c, r.g.Model(types.User{}).Count(&c).Error
}

func (r *UserRepository) GetUserByName(name string) (*types.User, error) {
	return FindSingle[types.User](r.g.Where(types.User{Name: name}))
}

func (r *UserRepository) GetUserByExternalIdAndType(external_id string, ut types.UserType) (*types.User, error) {
	return FindSingle[types.User](r.g.Where(types.User{ExternalID: &external_id, Type: ut}))
}

func (r *UserRepository) GetUsers(limit, offset int) ([]*types.User, error) {
	list := make([]*types.User, 0)
	err := r.g.Offset(offset).Limit(limit).Find(&list).Error
	return list, err
}

func (r *UserRepository) CreateUser(user *types.User) error {
	if user.UID == "" {
		user.UID = uuid.NewString()
	}
	return r.g.Create(user).Error
}

func (r *UserRepository) UpdateUser(user *types.User) error {
	return r.g.Updates(user).Error
}
