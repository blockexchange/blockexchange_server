package types

import (
	"fmt"
	"os"

	"github.com/luanti-golang/oauth"
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
	DefaultRole         UserRole
	// cache directory for the luanti wasm packs, a temp directory if empty
	WasmPackCacheDir string
	// cache directory for the git repositories of the wasm packs, "git" in WasmPackCacheDir if empty
	WasmGitCacheDir string
	// cache directory for the git repositories of the media sources, a temp directory if empty
	GitMediaCacheDir string
	// listen address for /metrics (for example ":9090"), served on the main listener if empty
	MetricsAddr string
}

func CreateConfig() *Config {
	cfg := &Config{
		Name:             os.Getenv("BLOCKEXCHANGE_NAME"),
		Owner:            os.Getenv("BLOCKEXCHANGE_OWNER"),
		Key:              os.Getenv("BLOCKEXCHANGE_KEY"),
		WebDev:           os.Getenv("WEBDEV") == "true",
		ExecuteJobs:      os.Getenv("EXECUTE_JOBS") == "true",
		BaseURL:          os.Getenv("BASE_URL"),
		CookiePath:       os.Getenv("BLOCKEXCHANGE_COOKIE_PATH"),
		CookieSecure:     os.Getenv("BLOCKEXCHANGE_COOKIE_SECURE") == "true",
		CookieName:       "blockexchange-v2",
		OauthLogin:       &OauthLogin{},
		DefaultRole:      UserRoleDefault,
		WasmPackCacheDir: os.Getenv("WASM_PACK_CACHE_DIR"),
		WasmGitCacheDir:  os.Getenv("WASM_GIT_CACHE_DIR"),
		GitMediaCacheDir: os.Getenv("GIT_MEDIA_CACHE_DIR"),
		MetricsAddr:      os.Getenv("METRICS_ADDR"),
	}

	if os.Getenv("DEFAULT_ROLE") != "" {
		cfg.DefaultRole = UserRole(os.Getenv("DEFAULT_ROLE"))
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

func (cfg *Config) Validate() error {
	switch cfg.DefaultRole {
	case UserRoleDefault, UserRoleAdmin:
	default:
		return fmt.Errorf("invalid DEFAULT_ROLE '%s', expected '%s' or '%s'", cfg.DefaultRole, UserRoleDefault, UserRoleAdmin)
	}
	return nil
}
