package api_test

import (
	"blockexchange/api"
	"blockexchange/core"
	"blockexchange/testutils"
	"blockexchange/types"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func NewTestApi(t *testing.T) *api.Api {
	db_ := testutils.CreateTestDatabase(t)
	cfg := types.CreateConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: "",
		DB:       0,
	})

	api, _, err := api.NewApi(db_, cfg, rdb)
	assert.NoError(t, err)
	return api
}

func Login(t *testing.T, r *http.Request, user *types.User) {
	permissions := []types.JWTPermission{
		types.JWTPermissionUpload,
		types.JWTPermissionManagement,
		types.JWTPermissionOverwrite,
		types.JWTPermissionAdmin,
	}
	c := core.New(types.CreateConfig(), nil)
	token, err := c.CreateJWT(user, permissions, time.Duration(1*time.Hour))
	assert.NoError(t, err)
	r.Header.Set("Authorization", token)
}
