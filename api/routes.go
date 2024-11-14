package api

import (
	"blockexchange/types"
	"net/http"

	"github.com/minetest-go/oauth"

	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
)

func (api *Api) SetupRoutes(r *mux.Router, cfg *types.Config) {

	// common api
	r.HandleFunc("/api/info", api.GetInfo).Methods(http.MethodGet)
	r.HandleFunc("/api/healthcheck", api.Healthcheck).Methods(http.MethodGet)

	// ui api
	r.HandleFunc("/api/login", api.DoLogin).Methods(http.MethodPost)
	r.HandleFunc("/api/login", api.GetLogin).Methods(http.MethodGet)
	r.HandleFunc("/api/login", api.Logout).Methods(http.MethodDelete)
	r.HandleFunc("/api/register", api.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/captcha", api.CreateCaptcha).Methods(http.MethodGet)
	r.PathPrefix("/api/captcha/").Handler(captcha.Server(300, 200))

	r.HandleFunc("/api/user-search", api.SearchUsers).Methods(http.MethodPost)
	r.HandleFunc("/api/user-count", api.CountUsers).Methods(http.MethodGet)
	r.HandleFunc("/api/user/{user_uid}/changepassword", api.Secure(api.ChangePassword)).Methods(http.MethodPost)
	r.HandleFunc("/api/user/{user_uid}/unlink-oauth", api.Secure(api.UnlinkOauth)).Methods(http.MethodPost)
	r.HandleFunc("/api/user/{user_uid}", api.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/api/user/{user_uid}/schemastars", api.CountUserSchemaStars).Methods(http.MethodGet)
	r.HandleFunc("/api/user/{user_uid}", api.Secure(api.SaveUser)).Methods(http.MethodPost)

	r.HandleFunc("/api/accesstoken", api.Secure(api.GetAccessTokens)).Methods(http.MethodGet)
	r.HandleFunc("/api/accesstoken", api.Secure(api.CreateAccessToken)).Methods(http.MethodPost)
	r.HandleFunc("/api/accesstoken/{id}", api.Secure(api.DeleteAccessToken)).Methods(http.MethodDelete)

	// mod uploads / ui downloads and manages

	r.HandleFunc("/api/media/mod/{modname}", api.Secure(api.CreateOrUpdateMod)).Methods(http.MethodPost)
	r.HandleFunc("/api/media/mod/{modname}", api.GetMod).Methods(http.MethodGet)
	r.HandleFunc("/api/media/mod", api.GetMods).Methods(http.MethodGet)
	r.HandleFunc("/api/media/mod/{modname}", api.Secure(api.DeleteMod)).Methods(http.MethodDelete)

	r.HandleFunc("/api/media/nodedef", api.Secure(api.CreateOrUpdateNodedefs)).Methods(http.MethodPost)
	r.HandleFunc("/api/media/nodedef", api.GetNodedefinitions).Methods(http.MethodGet)
	r.HandleFunc("/api/media/nodedef/{nodename}", api.GetNodedefinition).Methods(http.MethodGet)
	r.HandleFunc("/api/media/nodedef/{nodename}", api.Secure(api.DeleteNodedefinition)).Methods(http.MethodDelete)

	r.HandleFunc("/api/media/mediafile", api.Secure(api.CreateOrUpdateMediafiles)).Methods(http.MethodPost)
	r.HandleFunc("/api/media/mediafile/{name}", api.GetMediafile).Methods(http.MethodGet)
	r.HandleFunc("/api/media/mediafile/{name}", api.Secure(api.DeleteMediafile)).Methods(http.MethodDelete)

	// mod api
	r.HandleFunc("/api/token", api.RequestToken).Methods(http.MethodPost)

	r.HandleFunc("/api/export_we/{schema_uid}/{filename}", api.ExportWorldeditSchema).Methods(http.MethodGet)
	r.HandleFunc("/api/export_bx/{schema_uid}/{filename}", api.ExportBXSchema).Methods(http.MethodGet)
	r.HandleFunc("/api/import/{filename}", api.Secure(api.ImportSchematic)).Methods(http.MethodPost)

	r.HandleFunc("/api/tags", api.GetTags).Methods(http.MethodGet)

	r.HandleFunc("/api/schemamod/count", api.GetSchemaModCount).Methods(http.MethodGet)

	r.HandleFunc("/api/schema/{schema_uid}", api.GetSchema).Methods(http.MethodGet)
	r.HandleFunc("/api/schema/{schema_uid}", api.Secure(api.DeleteSchema)).Methods(http.MethodDelete)
	r.HandleFunc("/api/schema", api.Secure(api.CreateSchema)).Methods(http.MethodPost)
	r.HandleFunc("/api/schema/{schema_uid}", api.Secure(api.UpdateSchema)).Methods(http.MethodPut)
	r.HandleFunc("/api/schema/{schema_uid}/update", api.Secure(api.UpdateSchemaInfo)).Methods(http.MethodPost)
	r.HandleFunc("/api/schema/{schema_uid}/tags", api.Secure(api.UpdateSchemaTags)).Methods(http.MethodPost)

	r.HandleFunc("/api/collection", api.Secure(api.CreateOrUpdateCollection)).Methods(http.MethodPost)
	r.HandleFunc("/api/collection/by-username/{username}", api.GetCollectionsByUsername).Methods(http.MethodGet)
	r.HandleFunc("/api/collection/{collection_uid}", api.Secure(api.CreateOrUpdateCollection)).Methods(http.MethodPut)
	r.HandleFunc("/api/collection/{collection_uid}", api.Secure(api.DeleteCollection)).Methods(http.MethodDelete)
	r.HandleFunc("/api/collection/{collection_uid}", api.GetCollection).Methods(http.MethodGet)

	r.HandleFunc("/api/schema/{schema_uid}/mods", api.GetSchemaMods).Methods(http.MethodGet)
	r.HandleFunc("/api/schema/{schema_uid}/mods", api.Secure(api.CreateSchemaMods)).Methods(http.MethodPost)
	r.HandleFunc("/api/schema/{schema_uid}/mods/update", api.Secure(api.UpdateSchemaMods)).Methods(http.MethodPost)

	r.HandleFunc("/api/schema/{schema_uid}/star", api.Secure(api.GetSchemaStar)).Methods(http.MethodGet)
	r.HandleFunc("/api/schema/{schema_uid}/star/count", api.CountSchemaStars).Methods(http.MethodGet)
	r.HandleFunc("/api/schema/{schema_uid}/star", api.Secure(api.StarSchema)).Methods(http.MethodPut)
	r.HandleFunc("/api/schema/{schema_uid}/star", api.Secure(api.UnStarSchema)).Methods(http.MethodDelete)

	r.HandleFunc("/api/schema/{schema_uid}/screenshots", api.GetSchemaScreenshots)
	r.HandleFunc("/api/schema/{schema_uid}/screenshots/{screenshot_uid}", api.GetScreenshotByID)
	r.HandleFunc("/api/schema/{schema_uid}/screenshot/update", api.Secure(api.UpdateSchemaPreview)).Methods(http.MethodPost)
	r.HandleFunc("/api/schema/{schema_uid}/screenshot", api.GetFirstSchemaScreenshot)

	r.HandleFunc("/api/schema/{schema_uid}/pull", api.GetSchemaPull).Methods(http.MethodGet)
	r.HandleFunc("/api/schema/{schema_uid}/pull", api.Secure(api.CreateSchemaPull)).Methods(http.MethodPost)

	r.HandleFunc("/api/schema/{schema_uid}/pullclients", api.GetSchemaPullClients).Methods(http.MethodGet)

	r.HandleFunc("/api/search/schema", api.SearchSchema).Methods(http.MethodPost)
	r.HandleFunc("/api/count/schema", api.CountSchema).Methods(http.MethodPost)

	r.HandleFunc("/api/schemapart", api.Secure(api.CreateSchemaPart)).Methods(http.MethodPost)
	r.HandleFunc("/api/schemapart/{schema_uid}/{x}/{y}/{z}", api.GetSchemaPart).Methods(http.MethodGet)
	r.HandleFunc("/api/schemapart/{schema_uid}/{x}/{y}/{z}", api.Secure(api.DeleteSchemaPart)).Methods(http.MethodDelete)
	r.HandleFunc("/api/schemapart/{schema_uid}/{x}/{y}/{z}/delete", api.Secure(api.DeleteSchemaPart)).Methods(http.MethodPost)
	r.HandleFunc("/api/schemapart_chunk/{schema_uid}/{x}/{y}/{z}", api.GetSchemaPartChunk)
	r.HandleFunc("/api/schemapart_next/{schema_uid}/{x}/{y}/{z}", api.GetNextSchemaPart)
	r.HandleFunc("/api/schemapart_next/by-mtime/{schema_uid}/{mtime}", api.GetNextSchemaPartByMtime)
	r.HandleFunc("/api/schemapart_count/by-mtime/{schema_uid}/{mtime}", api.CountNextSchemaPartByMtime)
	r.HandleFunc("/api/schemapart_first/{schema_uid}", api.GetFirstSchemaPart)

	// oauth
	if cfg.GithubOAuthConfig != nil {
		oauth_handler := oauth.NewHandler(api.OauthCallback, cfg.GithubOAuthConfig)
		cfg.OauthLogin.Github = oauth_handler.LoginURL()
		r.Handle("/api/oauth_callback/github", oauth_handler)
	}

	if cfg.CDBOAuthConfig != nil {
		oauth_handler := oauth.NewHandler(api.OauthCallback, cfg.CDBOAuthConfig)
		cfg.OauthLogin.CDB = oauth_handler.LoginURL()
		r.Handle("/api/oauth_callback/cdb", oauth_handler)
	}

	if cfg.DiscordOAuthConfig != nil {
		oauth_handler := oauth.NewHandler(api.OauthCallback, cfg.DiscordOAuthConfig)
		cfg.OauthLogin.Discord = oauth_handler.LoginURL()
		r.Handle("/api/oauth_callback/discord", oauth_handler)
	}

	if cfg.MesehubOAuthConfig != nil {
		oauth_handler := oauth.NewHandler(api.OauthCallback, cfg.MesehubOAuthConfig)
		cfg.OauthLogin.Mesehub = oauth_handler.LoginURL()
		r.Handle("/api/oauth_callback/mesehub", oauth_handler)
	}

	if cfg.CodebergOAuthConfig != nil {
		oauth_handler := oauth.NewHandler(api.OauthCallback, cfg.CodebergOAuthConfig)
		cfg.OauthLogin.Codeberg = oauth_handler.LoginURL()
		r.Handle("/api/oauth_callback/codeberg", oauth_handler)
	}
}
