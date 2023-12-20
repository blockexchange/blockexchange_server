package types

import (
	"blockexchange/oauth"
	"os"
)

type Config struct {
	WebDev             bool
	BaseURL            string
	Name               string
	Owner              string
	Key                string
	GithubOAuthConfig  *oauth.OAuthConfig
	CDBOAuthConfig     *oauth.OAuthConfig
	DiscordOAuthConfig *oauth.OAuthConfig
	MesehubOAuthConfig *oauth.OAuthConfig
	CookiePath         string
	CookieSecure       bool
	CookieName         string
	RedisHost          string
	RedisPort          string
}

func CreateConfig() *Config {
	cfg := &Config{
		Name:         os.Getenv("BLOCKEXCHANGE_NAME"),
		Owner:        os.Getenv("BLOCKEXCHANGE_OWNER"),
		Key:          os.Getenv("BLOCKEXCHANGE_KEY"),
		WebDev:       os.Getenv("WEBDEV") == "true",
		BaseURL:      os.Getenv("BASE_URL"),
		CookiePath:   os.Getenv("BLOCKEXCHANGE_COOKIE_PATH"),
		CookieSecure: os.Getenv("BLOCKEXCHANGE_COOKIE_SECURE") == "true",
		CookieName:   "blockexchange",
		RedisHost:    os.Getenv("REDIS_HOST"),
		RedisPort:    os.Getenv("REDIS_PORT"),
	}

	if os.Getenv("DISCORD_APP_ID") != "" {
		cfg.DiscordOAuthConfig = &oauth.OAuthConfig{
			ClientID: os.Getenv("DISCORD_APP_ID"),
			Secret:   os.Getenv("DISCORD_APP_SECRET"),
		}
	}

	if os.Getenv("CDB_APP_ID") != "" {
		cfg.CDBOAuthConfig = &oauth.OAuthConfig{
			ClientID: os.Getenv("CDB_APP_ID"),
			Secret:   os.Getenv("CDB_APP_SECRET"),
		}
	}

	if os.Getenv("GITHUB_APP_ID") != "" {
		cfg.GithubOAuthConfig = &oauth.OAuthConfig{
			ClientID: os.Getenv("GITHUB_APP_ID"),
			Secret:   os.Getenv("GITHUB_APP_SECRET"),
		}
	}

	if os.Getenv("MESEHUB_APP_ID") != "" {
		cfg.MesehubOAuthConfig = &oauth.OAuthConfig{
			ClientID: os.Getenv("MESEHUB_APP_ID"),
			Secret:   os.Getenv("MESEHUB_APP_SECRET"),
		}
	}

	return cfg
}
