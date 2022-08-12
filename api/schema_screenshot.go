package api

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

func (api Api) GetFirstSchemaScreenshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_id, err := strconv.Atoi(vars["schema_id"])
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	screenshots, err := api.SchemaScreenshotRepo.GetBySchemaID(int64(schema_id))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if len(screenshots) == 0 {
		SendError(w, 404, "no screenshots found")
		return
	}

	screenshot := screenshots[0]
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

		cache_key := fmt.Sprintf("screenshot_%d_%d_%d", schema_id, height, width)
		data, err := api.Cache.Get(cache_key)
		if err != nil {
			// cache error
			SendJson(w, err.Error())
			return
		}

		if data != nil && r.URL.Query().Get("cache") != "false" {
			// cached data
			w.Header().Set("Cache-Control", "max-age=345600")
			w.Header().Set("Content-Type", "image/png")
			w.Write(data)
			return
		}

		img, err := png.Decode(bytes.NewReader(screenshot.Data))
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		newImage := resize.Resize(uint(width), uint(height), img, resize.Bilinear)

		buf := bytes.NewBuffer([]byte{})
		err = png.Encode(buf, newImage)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		// cache result
		err = api.Cache.Set(cache_key, buf.Bytes(), 10*time.Minute)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "max-age=345600")
		w.Write(buf.Bytes())
		return
	}

	w.Header().Set("Content-Type", screenshot.Type)
	w.Header().Set("Cache-Control", "max-age=345600")
	w.WriteHeader(200)
	w.Write(screenshot.Data)
}
