package web

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api Api) GetSchemaScreenshotByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	screenshot, err := api.SchemaScreenshotRepo.GetByID(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", screenshot.Type)
	w.WriteHeader(200)
	w.Write(screenshot.Data)
}
