package api

import (
	"blockexchange/types"
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

func getScreenshotPrefix(schema_uid string) string {
	return fmt.Sprintf("screenshot_%s_", schema_uid)
}

func createScreenshotKey(schema_uid string, height int, width int) string {
	return fmt.Sprintf("%s%d_%d", getScreenshotPrefix(schema_uid), height, width)
}

func (api Api) GetSchemaScreenshots(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	screenshots, err := api.SchemaScreenshotRepo.GetAllBySchemaUID(schema_uid)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	for _, s := range screenshots {
		s.Data = nil
	}

	Send(w, screenshots, nil)
}

func (api Api) GetScreenshotByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	screenshot_uid := vars["screenshot_uid"]
	screenshot, err := api.SchemaScreenshotRepo.GetByUID(screenshot_uid)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if screenshot == nil {
		SendError(w, 404, "no screenshot found")
		return
	}
	w.Header().Set("Content-Type", screenshot.Type)
	w.Header().Set("Cache-Control", "max-age=345600")
	w.WriteHeader(200)
	w.Write(screenshot.Data)
}

func (api Api) GetFirstSchemaScreenshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	screenshot, err := api.SchemaScreenshotRepo.GetLatestBySchemaUID(schema_uid)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if screenshot == nil {
		SendError(w, 404, "no screenshot found")
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

		cache_key := createScreenshotKey(schema_uid, height, width)
		data, err := api.Cache.Get(cache_key)
		if err != nil {
			// cache error
			SendJson(w, err.Error())
			return
		}

		if data != nil && r.URL.Query().Get("cache") != "false" {
			// cached data
			w.Header().Set("Cache-Control", "max-age=1800") // 30 minutes clientside cache
			w.Header().Set("Content-Type", "image/png")
			w.Write(data)
			return
		}

		img, err := png.Decode(bytes.NewReader(screenshot.Data))
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		oldWidth := img.Bounds().Max.X
		oldHeight := img.Bounds().Max.Y
		ratio := min(float64(width)/float64(oldWidth), float64(height)/float64(oldHeight))
		height = int(float64(oldHeight) * ratio)
		width = int(float64(oldWidth) * ratio)
		fmt.Printf("oldWidth: %d, oldHeight: %d, ratio: %.2f\n", oldWidth, oldHeight, ratio)

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
		//w.Header().Set("Cache-Control", "max-age=345600")
		w.Write(buf.Bytes())
		return
	}

	w.Header().Set("Content-Type", screenshot.Type)
	//w.Header().Set("Cache-Control", "max-age=345600")
	w.WriteHeader(200)
	w.Write(screenshot.Data)
}

func (api Api) UpdateSchemaPreview(w http.ResponseWriter, r *http.Request, ctx *SecureContext) {
	vars := mux.Vars(r)
	schema_uid := vars["schema_uid"]

	schema, err := api.SchemaRepo.GetSchemaByUID(schema_uid)
	if err != nil {
		SendError(w, 500, "GetSchemaById::"+err.Error())
		return
	}

	if schema.UserUID != ctx.Claims.UserUID && !ctx.Claims.HasPermission(types.JWTPermissionAdmin) {
		SendError(w, 403, "you are not the owner of the schema")
		return
	}

	// update screenshot
	_, err = api.core.UpdatePreview(schema)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// clear cache
	err = api.Cache.RemovePrefix(getScreenshotPrefix(schema_uid))
	Send(w, true, err)
}
