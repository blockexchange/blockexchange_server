package types_test

import (
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDefaultRole(t *testing.T) {
	t.Setenv("DEFAULT_ROLE", "")
	cfg := types.CreateConfig()
	assert.Equal(t, types.UserRoleDefault, cfg.DefaultRole)
	assert.NoError(t, cfg.Validate())

	t.Setenv("DEFAULT_ROLE", "ADMIN")
	cfg = types.CreateConfig()
	assert.Equal(t, types.UserRoleAdmin, cfg.DefaultRole)
	assert.NoError(t, cfg.Validate())

	t.Setenv("DEFAULT_ROLE", "ADMN")
	cfg = types.CreateConfig()
	assert.Error(t, cfg.Validate())
}
