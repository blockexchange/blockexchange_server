package api

import (
	"blockexchange/core"
	"blockexchange/schematic/worldedit"
	"blockexchange/types"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (api *Api) incrementDownloadstats(schema_id int64, r *http.Request) error {
	origin := r.Header.Get("X-Forwarded-For")
	origin_parts := strings.Split(origin, ",")
	if len(origin_parts) >= 1 {
		ip := strings.TrimSpace(origin_parts[0])
		cache_key := fmt.Sprintf("download_marker/%d/%s", schema_id, ip)

		m, err := api.Cache.Get(cache_key)
		if err != nil {
			return err
		}
		if m != nil {
			//already incremented stat for this ip, skip
			return nil
		}
		// mark as incremented
		err = api.Cache.Set(cache_key, []byte{0x00}, time.Hour*24)
		if err != nil {
			return err
		}
	}

	err := api.SchemaRepo.IncrementDownloads(schema_id)
	if err != nil {
		return err
	}
	return nil
}

func (api *Api) ExportWorldeditSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", vars["filename"]))
	e := worldedit.NewExporter(w)

	schemapart, err := api.SchemaPartRepo.GetFirstBySchemaID(int64(id))
	err = e.Export(schemapart)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	for {
		schemapart, err = api.SchemaPartRepo.GetNextBySchemaIDAndOffset(int64(id), schemapart.OffsetX, schemapart.OffsetY, schemapart.OffsetZ)
		if schemapart != nil {
			err = e.Export(schemapart)
			if err != nil {
				SendError(w, 500, err.Error())
				return
			}
		} else {
			break
		}
	}

	err = e.Close()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	err = api.incrementDownloadstats(int64(id), r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}

func (api *Api) ExportBXSchema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	schema, err := api.SchemaRepo.GetSchemaById(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if schema == nil {
		SendError(w, 404, "not found")
		return
	}

	schemamods, err := api.SchemaModRepo.GetSchemaModsBySchemaID(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	var schemapart *types.SchemaPart
	it := func() (*types.SchemaPart, error) {
		var err error
		if schemapart == nil {
			schemapart, err = api.SchemaPartRepo.GetFirstBySchemaID(int64(id))
		} else {
			schemapart, err = api.SchemaPartRepo.GetNextBySchemaIDAndOffset(int64(id), schemapart.OffsetX, schemapart.OffsetY, schemapart.OffsetZ)
		}
		return schemapart, err
	}

	err = core.ExportBXSchema(w, schema, schemamods, it)
	if err != nil {
		SendError(w, 500, err.Error())
	}

	err = api.incrementDownloadstats(int64(id), r)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}
