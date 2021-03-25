package web

import (
	"blockexchange/core"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct{}

func (r MockUserRepo) GetUserById(id int64) (*types.User, error) {
	return nil, nil
}
func (r MockUserRepo) GetUserByName(name string) (*types.User, error) {
	var user *types.User

	if name == "user" {
		hash, err := bcrypt.GenerateFromPassword([]byte("pw"), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user = &types.User{
			ID:   0,
			Name: name,
			Hash: string(hash),
			Type: types.UserTypeLocal,
		}
	}
	return user, nil
}
func (r MockUserRepo) GetUserByExternalId(external_id string) (*types.User, error) {
	return nil, nil
}
func (r MockUserRepo) GetUsers() ([]types.User, error) {
	return make([]types.User, 0), nil
}
func (r MockUserRepo) CreateUser(user *types.User) error {
	return nil
}
func (r MockUserRepo) UpdateUser(user *types.User) error {
	return nil
}

func createTestApi() *Api {
	return &Api{
		UserRepo: &MockUserRepo{},
	}
}

func TestValidUserLogin(t *testing.T) {
	api := createTestApi()

	login := types.Login{}
	login.Username = "user"
	login.Password = "pw"
	data, err := json.Marshal(login)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	r := httptest.NewRequest("GET", "http://", bytes.NewReader(data))
	w := httptest.NewRecorder()

	api.PostLogin(w, r)

	assert.NotNil(t, w.Body)
	info, err := core.ParseJWT(w.Body.String())
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "user", info.Username)
}

func TestInvalidUserLogin(t *testing.T) {
	api := createTestApi()

	login := types.Login{}
	login.Username = "user"
	login.Password = "pw2"
	data, err := json.Marshal(login)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	r := httptest.NewRequest("GET", "http://", bytes.NewReader(data))
	w := httptest.NewRecorder()

	api.PostLogin(w, r)

	assert.NotNil(t, w.Body)
	info, err := core.ParseJWT(w.Body.String())
	assert.Error(t, err)
	assert.Nil(t, info)
}
