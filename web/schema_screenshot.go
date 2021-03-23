package web

import (
	"bytes"
	"image/png"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"github.com/sirupsen/logrus"
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

	if screenshot == nil {
		w.WriteHeader(404)
		return
	}

	if r.URL.Query().Get("height") != "" {
		height, err := strconv.Atoi(r.URL.Query().Get("height"))
		if err != nil || height < 0 || height > 2048 {
			SendError(w, 500, err.Error())
			return
		}

		width, err := strconv.Atoi(r.URL.Query().Get("width"))
		if err != nil || width < 0 || width > 2048 {
			SendError(w, 500, err.Error())
			return
		}

		logrus.WithFields(logrus.Fields{
			"width":  width,
			"height": height,
		}).Trace("api::GetSchemaScreenshotByID rescaling image")

		img, err := png.Decode(bytes.NewReader(screenshot.Data))
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		newImage := resize.Resize(uint(width), uint(height), img, resize.Bilinear)
		w.Header().Set("Content-Type", "image/png")
		err = png.Encode(w, newImage)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		return
	}

	w.Header().Set("Content-Type", screenshot.Type)
	w.WriteHeader(200)
	w.Write(screenshot.Data)
}

func (api Api) GetSchemaScreenshots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	screenshots, err := api.SchemaScreenshotRepo.GetBySchemaID(int64(id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	for i := range screenshots {
		screenshots[i].Data = nil
	}

	SendJson(w, screenshots)
}
