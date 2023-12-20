package api

import (
	"blockexchange/oauth"
	"blockexchange/types"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
)

func (api *Api) SetupRoutes(r *mux.Router, cfg *types.Config) {

	// common api
	r.Handle("/api/info", InfoHandler{Config: cfg})
	r.HandleFunc("/api/healthcheck", api.Healthcheck)

	// ui api
	r.HandleFunc("/api/login", api.DoLogin).Methods(http.MethodPost)
	r.HandleFunc("/api/login", api.GetLogin).Methods(http.MethodGet)
	r.HandleFunc("/api/login", api.Logout).Methods(http.MethodDelete)
	r.HandleFunc("/api/register", api.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/captcha", api.CreateCaptcha).Methods(http.MethodGet)
	r.PathPrefix("/api/captcha/").Handler(captcha.Server(300, 200))

	r.HandleFunc("/api/user-search", api.SearchUsers).Methods(http.MethodPost)
	r.HandleFunc("/api/user-count", api.CountUsers).Methods(http.MethodGet)
	r.HandleFunc("/api/user/{user_id}", api.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/api/user/{user_id}", api.Secure(api.SaveUser)).Methods(http.MethodPost)

	r.HandleFunc("/api/accesstoken", api.Secure(api.GetAccessTokens)).Methods(http.MethodGet)
	r.HandleFunc("/api/accesstoken", api.Secure(api.CreateAccessToken)).Methods(http.MethodPost)
	r.HandleFunc("/api/accesstoken/{id}", api.Secure(api.DeleteAccessToken)).Methods(http.MethodDelete)

	// mod api
	r.HandleFunc("/api/token", api.RequestToken).Methods(http.MethodPost)

	r.HandleFunc("/api/export_we/{id}/{filename}", api.ExportWorldeditSchema).Methods(http.MethodGet)
	r.HandleFunc("/api/export_bx/{id}/{filename}", api.ExportBXSchema).Methods(http.MethodGet)
	r.HandleFunc("/api/import/{filename}", api.Secure(api.ImportSchematic)).Methods(http.MethodPost)

	r.HandleFunc("/api/tags", api.GetTags).Methods(http.MethodGet)

	r.HandleFunc("/api/schema/{id}", api.GetSchema).Methods(http.MethodGet)
	r.HandleFunc("/api/schema/{id}", api.Secure(api.DeleteSchema)).Methods(http.MethodDelete)
	r.HandleFunc("/api/schema", api.Secure(api.CreateSchema)).Methods(http.MethodPost)
	r.HandleFunc("/api/schema/{id}", api.Secure(api.UpdateSchema)).Methods(http.MethodPut)
	r.HandleFunc("/api/schema/{id}/mods", api.GetSchemaMods).Methods(http.MethodGet)
	r.HandleFunc("/api/schema/{id}/mods", api.Secure(api.CreateSchemaMods)).Methods(http.MethodPost)
	r.HandleFunc("/api/schema/{id}/update", api.Secure(api.UpdateSchemaInfo)).Methods(http.MethodPost)
	r.HandleFunc("/api/schema/{id}/tags", api.Secure(api.UpdateSchemaTags)).Methods(http.MethodPost)

	r.HandleFunc("/api/schema/{schema_id}/star", api.Secure(api.GetSchemaStar)).Methods(http.MethodGet)
	r.HandleFunc("/api/schema/{schema_id}/star/count", api.CountSchemaStars).Methods(http.MethodGet)
	r.HandleFunc("/api/schema/{schema_id}/star", api.Secure(api.StarSchema)).Methods(http.MethodPut)
	r.HandleFunc("/api/schema/{schema_id}/star", api.Secure(api.UnStarSchema)).Methods(http.MethodDelete)

	r.HandleFunc("/api/schema/{schema_id}/screenshot/update", api.Secure(api.UpdateSchemaPreview)).Methods(http.MethodPost)
	r.HandleFunc("/api/schema/{schema_id}/screenshot", api.GetFirstSchemaScreenshot)

	r.HandleFunc("/api/search/schema/byname/{user_name}/{schema_name}", api.SearchSchemaByNameAndUser)
	r.HandleFunc("/api/search/schema", api.SearchSchema).Methods(http.MethodPost)
	r.HandleFunc("/api/count/schema", api.CountSchema).Methods(http.MethodPost)

	r.HandleFunc("/api/schemapart", api.Secure(api.CreateSchemaPart)).Methods(http.MethodPost)
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}", api.GetSchemaPart).Methods(http.MethodGet)
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}", api.Secure(api.DeleteSchemaPart)).Methods(http.MethodDelete)
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}/delete", api.Secure(api.DeleteSchemaPart)).Methods(http.MethodPost)
	r.HandleFunc("/api/schemapart_chunk/{schema_id}/{x}/{y}/{z}", api.GetSchemaPartChunk)
	r.HandleFunc("/api/schemapart_next/{schema_id}/{x}/{y}/{z}", api.GetNextSchemaPart)
	r.HandleFunc("/api/schemapart_next/by-mtime/{schema_id}/{mtime}", api.GetNextSchemaPartByMtime)
	r.HandleFunc("/api/schemapart_count/by-mtime/{schema_id}/{mtime}", api.CountNextSchemaPartByMtime)
	r.HandleFunc("/api/schemapart_first/{schema_id}", api.GetFirstSchemaPart)

	// oauth
	if cfg.GithubOAuthConfig != nil {
		oauth_handler := &oauth.OauthHandler{
			Impl:     &oauth.GithubOauth{},
			Config:   cfg.GithubOAuthConfig,
			BaseURL:  api.cfg.BaseURL,
			Callback: api.OauthCallback,
		}
		r.Handle("/api/oauth_callback/github", oauth_handler)
	}

	if cfg.CDBOAuthConfig != nil {
		oauth_handler := &oauth.OauthHandler{
			Impl:     &oauth.CDBOauth{},
			Config:   cfg.CDBOAuthConfig,
			BaseURL:  api.cfg.BaseURL,
			Callback: api.OauthCallback,
		}
		r.Handle("/api/oauth_callback/cdb", oauth_handler)
	}

	if cfg.DiscordOAuthConfig != nil {
		oauth_handler := &oauth.OauthHandler{
			Impl:     &oauth.DiscordOauth{},
			Config:   cfg.DiscordOAuthConfig,
			BaseURL:  api.cfg.BaseURL,
			Callback: api.OauthCallback,
		}
		r.Handle("/api/oauth_callback/discord", oauth_handler)
	}

	if cfg.MesehubOAuthConfig != nil {
		oauth_handler := &oauth.OauthHandler{
			Impl:     &oauth.MesehubOauth{},
			Config:   cfg.MesehubOAuthConfig,
			BaseURL:  api.cfg.BaseURL,
			Callback: api.OauthCallback,
		}
		r.Handle("/api/oauth_callback/mesehub", oauth_handler)
	}
}
