package testutils

import (
	"blockexchange/core"
	"blockexchange/db"
	"blockexchange/types"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func CreateName(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CreateUser(repo db.UserRepository, t *testing.T, user *types.User) *types.User {
	if user == nil {
		// create new user
		user = &types.User{}
	}
	// set defaults
	if user.Name == "" {
		user.Name = CreateName(10)
	}
	if user.Type == "" {
		user.Type = types.UserTypeLocal
	}

	assert.NoError(t, repo.CreateUser(user))

	return user
}

func CreateSchema(repo db.SchemaRepository, t *testing.T, user *types.User, schema *types.Schema) *types.Schema {
	if schema == nil {
		schema = &types.Schema{
			Name: CreateName(10),
		}
	}

	schema.UserID = user.ID
	assert.NoError(t, repo.CreateSchema(schema))
	return schema
}

func CreateSchemaScreenshot(repo db.SchemaScreenshotRepository, t *testing.T, schema *types.Schema, screenshot *types.SchemaScreenshot) *types.SchemaScreenshot {
	if screenshot == nil {
		screenshot = &types.SchemaScreenshot{}
	}

	screenshot.SchemaID = schema.ID
	assert.NoError(t, repo.Create(screenshot))
	return screenshot
}

func CreateAccessToken(repo db.AccessTokenRepository, t *testing.T, token *types.AccessToken) *types.AccessToken {
	if token == nil {
		token = &types.AccessToken{}
	}

	assert.NoError(t, repo.CreateAccessToken(token))
	return token
}

func Login(t *testing.T, r *http.Request, user *types.User) {
	permissions := []types.JWTPermission{
		types.JWTPermissionUpload,
		types.JWTPermissionManagement,
		types.JWTPermissionOverwrite,
		types.JWTPermissionAdmin,
	}
	token, err := core.CreateJWT(user, permissions, (time.Now().Unix()+3600)+1000)
	assert.NoError(t, err)
	r.Header.Set("Authorization", token)
}
