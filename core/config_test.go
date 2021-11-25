package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateConfig(t *testing.T) {
	cfg, err := CreateConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	os.Setenv("DISCORD_APP_ID", "dummy")
	os.Setenv("DISCORD_APP_SECRET", "xy")

	cfg, err = CreateConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.DiscordOAuthConfig)
	assert.Equal(t, "dummy", cfg.DiscordOAuthConfig.ClientID)
	assert.Equal(t, "xy", cfg.DiscordOAuthConfig.Secret)
}
