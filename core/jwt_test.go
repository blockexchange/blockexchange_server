package core_test

import (
	"blockexchange/core"
	"blockexchange/types"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func contains(s []types.JWTPermission, e types.JWTPermission) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestGetPermissions(t *testing.T) {
	permissions := core.GetPermissions(&types.User{
		Role: types.UserRoleDefault,
	}, true)

	assert.True(t, contains(permissions, types.JWTPermissionManagement))
	assert.True(t, contains(permissions, types.JWTPermissionUpload))
	assert.True(t, contains(permissions, types.JWTPermissionOverwrite))
	assert.False(t, contains(permissions, types.JWTPermissionAdmin))

	permissions = core.GetPermissions(&types.User{
		Role: types.UserRoleDefault,
	}, false)

	assert.False(t, contains(permissions, types.JWTPermissionManagement))
	assert.True(t, contains(permissions, types.JWTPermissionUpload))
	assert.True(t, contains(permissions, types.JWTPermissionOverwrite))
	assert.False(t, contains(permissions, types.JWTPermissionAdmin))

	permissions = core.GetPermissions(&types.User{
		Role: types.UserRoleAdmin,
	}, true)

	assert.True(t, contains(permissions, types.JWTPermissionManagement))
	assert.True(t, contains(permissions, types.JWTPermissionUpload))
	assert.True(t, contains(permissions, types.JWTPermissionOverwrite))
	assert.True(t, contains(permissions, types.JWTPermissionAdmin))

	permissions = core.GetPermissions(&types.User{
		Role: types.UserRoleAdmin,
	}, false)

	assert.False(t, contains(permissions, types.JWTPermissionManagement))
	assert.True(t, contains(permissions, types.JWTPermissionUpload))
	assert.True(t, contains(permissions, types.JWTPermissionOverwrite))
	assert.False(t, contains(permissions, types.JWTPermissionAdmin))
}

func TestCreateJWT(t *testing.T) {
	uid := uuid.NewString()
	user := types.User{
		Name: "dummy",
		UID:  uid,
		Type: types.UserTypeLocal,
	}
	permissions := []types.JWTPermission{types.JWTPermissionUpload}
	c := core.New(types.CreateConfig(), nil)

	token, err := c.CreateJWT(&user, permissions, time.Duration(30*time.Second))
	if err != nil {
		t.Fatal(err)
	}

	info, err := c.ParseJWT(token)
	if err != nil {
		t.Fatal(err)
	}

	if info.Username != "dummy" {
		t.Fatal("username mismatch")
	}

	//TODO: check permissions, etc
}
