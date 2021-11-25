package web

import (
	"embed"
	"net/http"
	"os"

	"blockexchange/core"
	"blockexchange/public"
	"blockexchange/web/oauth"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(db_ *sqlx.DB) {

	r := mux.NewRouter()

	// cache
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	var cache core.Cache
	if redis_host != "" && redis_port != "" {
		cache = core.NewRedisCache(redis_host + ":" + redis_port)
	} else {
		cache = core.NewNoOpCache()
	}

	api := NewApi(db_, cache)
	SetupRoutes(r, api)

	handler := cors.Default().Handler(r)
	http.Handle("/", handler)

	// metrics
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func SetupRoutes(r *mux.Router, api *Api) {

	github_oauth := oauth.NewHandler(&oauth.GithubOauth{}, api.UserRepo, api.AccessTokenRepo)
	discord_oauth := oauth.NewHandler(&oauth.DiscordOauth{}, api.UserRepo, api.AccessTokenRepo)
	mesehub_oauth := oauth.NewHandler(&oauth.MesehubOauth{}, api.UserRepo, api.AccessTokenRepo)

	// api surface
	r.HandleFunc("/api/info", InfoEndpoint)
	r.HandleFunc("/api/register", api.Register).Methods("POST")
	r.HandleFunc("/api/token", api.PostLogin).Methods("POST")
	r.HandleFunc("/api/oauth_callback/github", github_oauth.Handle)
	r.HandleFunc("/api/oauth_callback/discord", discord_oauth.Handle)
	r.HandleFunc("/api/oauth_callback/mesehub", mesehub_oauth.Handle)

	r.HandleFunc("/api/validate_username", api.PostValidateUsername).Methods("POST")
	r.HandleFunc("/api/user", api.GetUsers).Methods("GET")
	r.HandleFunc("/api/user/{id}", Secure(api.UpdateUser)).Methods("POST")

	r.HandleFunc("/api/export_we/{username}/{schemaname}/{filename}", api.ExportWorldeditSchema).Methods("GET")
	r.HandleFunc("/api/export_bx/{username}/{schemaname}/{filename}", api.ExportBXSchema).Methods("GET")

	r.HandleFunc("/api/schema/{username}/{schemaname}", api.GetSchema).Methods("GET")
	r.HandleFunc("/api/schema", Secure(api.CreateSchema)).Methods("POST")
	r.HandleFunc("/api/schema", Secure(api.UpdateSchema)).Methods("PUT")
	r.HandleFunc("/api/schema/{username}/{schemaname}", Secure(api.DeleteSchema)).Methods("DELETE")
	r.HandleFunc("/api/schema/{username}/{schemaname}/update", Secure(api.UpdateSchemaInfo)).Methods("POST")

	r.HandleFunc("/api/schema/{id}/mods", api.GetSchemaMods).Methods("GET")
	r.HandleFunc("/api/schema/{id}/mods", Secure(api.CreateSchemaMods)).Methods("POST")

	r.HandleFunc("/api/schema/{schema_id}/star", Secure(api.CreateSchemaStar)).Methods("POST")
	r.HandleFunc("/api/schema/{schema_id}/star", Secure(api.DeleteSchemaStar)).Methods("DELETE")
	r.HandleFunc("/api/schema/{schema_id}/star", api.GetSchemaStars).Methods("GET")

	r.HandleFunc("/api/schema/{schema_id}/tag/{tag_id}", Secure(api.CreateSchemaTag)).Methods("PUT")
	r.HandleFunc("/api/schema/{schema_id}/tag/{tag_id}", Secure(api.DeleteSchemaTag)).Methods("DELETE")
	r.HandleFunc("/api/schema/{schema_id}/tag", api.GetSchemaTags).Methods("GET")

	r.HandleFunc("/api/collection/by-user_id/{user_id}", api.GetCollectionsByUserID).Methods("GET")
	r.HandleFunc("/api/collection", Secure(api.CreateCollection)).Methods("POST")
	r.HandleFunc("/api/collection/{id}", Secure(api.UpdateCollection)).Methods("PUT")
	r.HandleFunc("/api/collection/{id}", Secure(api.DeleteCollection)).Methods("DELETE")

	r.HandleFunc("/api/tag", api.GetTags).Methods("GET")
	r.HandleFunc("/api/tag", Secure(api.CreateTag)).Methods("POST")
	r.HandleFunc("/api/tag", Secure(api.UpdateTag)).Methods("PUT")
	r.HandleFunc("/api/tag/{id}", Secure(api.DeleteTag)).Methods("DELETE")

	r.HandleFunc("/api/schema/{schema_id}/screenshot/{id}", api.GetSchemaScreenshotByID)
	r.HandleFunc("/api/schema/{schema_id}/screenshot", api.GetSchemaScreenshots)

	r.HandleFunc("/api/static/schema/{user_name}/{schema_name}", api.GetStaticView)

	r.HandleFunc("/api/search/schema/byname/{user_name}/{schema_name}", api.SearchSchemaByNameAndUser)
	r.HandleFunc("/api/searchschema", api.SearchSchema).Methods("POST")
	r.HandleFunc("/api/searchrecent/{count}", api.SearchRecentSchemas)

	r.HandleFunc("/api/schemapart", Secure(api.CreateSchemaPart)).Methods("POST")
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}", api.GetSchemaPart)
	r.HandleFunc("/api/schemapart_chunk/{schema_id}/{x}/{y}/{z}", api.GetSchemaPartChunk)
	r.HandleFunc("/api/schemapart_next/{schema_id}/{x}/{y}/{z}", api.GetNextSchemaPart)
	r.HandleFunc("/api/schemapart_next/by-mtime/{schema_id}/{mtime}", api.GetNextSchemaPartByMtime)
	r.HandleFunc("/api/schemapart_first/{schema_id}", api.GetFirstSchemaPart)

	r.HandleFunc("/api/access_token", Secure(api.GetAccessTokens)).Methods("GET")
	r.HandleFunc("/api/access_token", Secure(api.PostAccessToken)).Methods("POST")
	r.HandleFunc("/api/access_token/{id}", Secure(api.DeleteAccessToken)).Methods("DELETE")

	// webdev flag
	useLocalfs := os.Getenv("WEBDEV") == "true"
	// static files
	r.PathPrefix("/").Handler(http.FileServer(getFileSystem(useLocalfs, public.Webapp)))
}

func getFileSystem(useLocalfs bool, content embed.FS) http.FileSystem {
	if useLocalfs {
		logrus.Print("using live mode")
		return http.FS(os.DirFS("public"))
	}

	logrus.Print("using embed mode")
	return http.FS(content)
}
