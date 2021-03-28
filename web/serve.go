package web

import (
	"embed"
	"log"
	"net/http"
	"os"

	"blockexchange/core"
	"blockexchange/public"
	"blockexchange/web/oauth"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(db_ *sqlx.DB) {

	// webdev flag
	useLocalfs := os.Getenv("WEBDEV") == "true"
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

	r.HandleFunc("/api/schema/{id}", api.GetSchema).Methods("GET")
	r.HandleFunc("/api/schema", Secure(api.CreateSchema)).Methods("POST")
	r.HandleFunc("/api/schema/{id}", Secure(api.UpdateSchema)).Methods("PUT")
	r.HandleFunc("/api/schema/{id}", Secure(api.DeleteSchema)).Methods("DELETE")
	r.HandleFunc("/api/schema/{id}/mods", api.GetSchemaMods).Methods("GET")
	r.HandleFunc("/api/schema/{id}/mods", Secure(api.CreateSchemaMods)).Methods("POST")
	r.HandleFunc("/api/schema/{id}/update", Secure(api.UpdateSchemaInfo)).Methods("POST")

	r.HandleFunc("/api/schema/{schema_id}/tag/{tag_id}", Secure(api.CreateSchemaTag)).Methods("PUT")
	r.HandleFunc("/api/schema/{schema_id}/tag/{tag_id}", Secure(api.DeleteSchemaTag)).Methods("DELETE")

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

	r.HandleFunc("/api/search/schema/byname/{user_name}/{schema_name}", api.SearchSchemaByNameAndUser)
	r.HandleFunc("/api/searchschema", api.SearchSchema).Methods("POST")
	r.HandleFunc("/api/searchrecent/{count}", api.SearchRecentSchemas)

	r.HandleFunc("/api/schemapart", Secure(api.CreateSchemaPart)).Methods("POST")
	r.HandleFunc("/api/schemapart/{schema_id}/{x}/{y}/{z}", api.GetSchemaPart)
	r.HandleFunc("/api/schemapart_next/{schema_id}/{x}/{y}/{z}", api.GetNextSchemaPart)
	r.HandleFunc("/api/schemapart_first/{schema_id}", api.GetFirstSchemaPart)

	r.HandleFunc("/api/access_token", Secure(api.GetAccessTokens)).Methods("GET")
	r.HandleFunc("/api/access_token", Secure(api.PostAccessToken)).Methods("POST")
	r.HandleFunc("/api/access_token/{id}", Secure(api.DeleteAccessToken)).Methods("DELETE")

	// static files
	r.PathPrefix("/").Handler(http.FileServer(getFileSystem(useLocalfs, public.Webapp)))
	http.Handle("/", r)

	// metrics
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func getFileSystem(useLocalfs bool, content embed.FS) http.FileSystem {
	if useLocalfs {
		log.Print("using live mode")
		return http.FS(os.DirFS("public"))
	}

	log.Print("using embed mode")
	return http.FS(content)
}
