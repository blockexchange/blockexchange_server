package types

import (
	"fmt"
	"os"

	"github.com/minetest-go/oauth"
)

type OauthLogin struct {
	Github   string `json:"github"`
	Discord  string `json:"discord"`
	Mesehub  string `json:"mesehub"`
	CDB      string `json:"cdb"`
	Codeberg string `json:"codeberg"`
}

type Config struct {
	WebDev              bool
	ExecuteJobs         bool
	BaseURL             string
	Name                string
	Owner               string
	Key                 string
	GithubOAuthConfig   *oauth.OAuthConfig
	OauthLogin          *OauthLogin
	CDBOAuthConfig      *oauth.OAuthConfig
	DiscordOAuthConfig  *oauth.OAuthConfig
	MesehubOAuthConfig  *oauth.OAuthConfig
	CodebergOAuthConfig *oauth.OAuthConfig
	CookiePath          string
	CookieSecure        bool
	CookieName          string
	RedisHost           string
	RedisPort           string
}

func CreateConfig() *Config {
	cfg := &Config{
		Name:         os.Getenv("BLOCKEXCHANGE_NAME"),
		Owner:        os.Getenv("BLOCKEXCHANGE_OWNER"),
		Key:          os.Getenv("BLOCKEXCHANGE_KEY"),
		WebDev:       os.Getenv("WEBDEV") == "true",
		ExecuteJobs:  os.Getenv("EXECUTE_JOBS") == "true",
		BaseURL:      os.Getenv("BASE_URL"),
		CookiePath:   os.Getenv("BLOCKEXCHANGE_COOKIE_PATH"),
		CookieSecure: os.Getenv("BLOCKEXCHANGE_COOKIE_SECURE") == "true",
		CookieName:   "blockexchange",
		RedisHost:    os.Getenv("REDIS_HOST"),
		RedisPort:    os.Getenv("REDIS_PORT"),
		OauthLogin:   &OauthLogin{},
	}

	if os.Getenv("DISCORD_APP_ID") != "" {
		cfg.DiscordOAuthConfig = &oauth.OAuthConfig{
			Provider:    oauth.ProviderTypeDiscord,
			ClientID:    os.Getenv("DISCORD_APP_ID"),
			Secret:      os.Getenv("DISCORD_APP_SECRET"),
			CallbackURL: fmt.Sprintf("%s/api/oauth_callback/discord", cfg.BaseURL),
		}
	}

	if os.Getenv("CDB_APP_ID") != "" {
		cfg.CDBOAuthConfig = &oauth.OAuthConfig{
			Provider:    oauth.ProviderTypeCDB,
			ClientID:    os.Getenv("CDB_APP_ID"),
			Secret:      os.Getenv("CDB_APP_SECRET"),
			CallbackURL: fmt.Sprintf("%s/api/oauth_callback/cdb", cfg.BaseURL),
		}
	}

	if os.Getenv("GITHUB_APP_ID") != "" {
		cfg.GithubOAuthConfig = &oauth.OAuthConfig{
			Provider:    oauth.ProviderTypeGithub,
			ClientID:    os.Getenv("GITHUB_APP_ID"),
			Secret:      os.Getenv("GITHUB_APP_SECRET"),
			CallbackURL: fmt.Sprintf("%s/api/oauth_callback/github", cfg.BaseURL),
		}
	}

	if os.Getenv("MESEHUB_APP_ID") != "" {
		cfg.MesehubOAuthConfig = &oauth.OAuthConfig{
			Provider:    oauth.ProviderTypeMesehub,
			ClientID:    os.Getenv("MESEHUB_APP_ID"),
			Secret:      os.Getenv("MESEHUB_APP_SECRET"),
			CallbackURL: fmt.Sprintf("%s/api/oauth_callback/mesehub", cfg.BaseURL),
		}
	}

	if os.Getenv("CODEBERG_APP_ID") != "" {
		cfg.CodebergOAuthConfig = &oauth.OAuthConfig{
			Provider:    oauth.ProviderTypeCodeberg,
			ClientID:    os.Getenv("CODEBERG_APP_ID"),
			Secret:      os.Getenv("CODEBERG_APP_SECRET"),
			CallbackURL: fmt.Sprintf("%s/api/oauth_callback/codeberg", cfg.BaseURL),
		}
	}

	return cfg
}
