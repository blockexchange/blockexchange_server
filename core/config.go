package core

import "os"

type OAuthConfig struct {
	ClientID string
	Secret   string
}

type Config struct {
	WebDev             bool
	EnableSignup       bool
	BaseURL            string
	Name               string
	Owner              string
	Key                string
	GithubOAuthConfig  *OAuthConfig
	DiscordOAuthConfig *OAuthConfig
	MesehubOAuthConfig *OAuthConfig
}

func CreateConfig() (*Config, error) {
	cfg := &Config{
		Name:         os.Getenv("BLOCKEXCHANGE_NAME"),
		Owner:        os.Getenv("BLOCKEXCHANGE_OWNER"),
		WebDev:       os.Getenv("WEBDEV") == "true",
		EnableSignup: os.Getenv("ENABLE_SIGNUP") == "true",
		BaseURL:      os.Getenv("BASE_URL"),
	}

	if os.Getenv("DISCORD_APP_ID") != "" {
		cfg.DiscordOAuthConfig = &OAuthConfig{
			ClientID: os.Getenv("DISCORD_APP_ID"),
			Secret:   os.Getenv("DISCORD_APP_SECRET"),
		}
	}

	if os.Getenv("GITHUB_APP_ID") != "" {
		cfg.GithubOAuthConfig = &OAuthConfig{
			ClientID: os.Getenv("GITHUB_APP_ID"),
			Secret:   os.Getenv("GITHUB_APP_SECRET"),
		}
	}

	if os.Getenv("MESEHUB_APP_ID") != "" {
		cfg.MesehubOAuthConfig = &OAuthConfig{
			ClientID: os.Getenv("MESEHUB_APP_ID"),
			Secret:   os.Getenv("MESEHUB_APP_SECRET"),
		}
	}

	return cfg, nil
}
