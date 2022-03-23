package web

import (
	"blockexchange/colormapping"
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

var cachedColormapping []byte

func (api *Api) GetColorMapping(w http.ResponseWriter, r *http.Request) {
	if cachedColormapping == nil {
		cm := colormapping.NewColorMapping()
		err := cm.LoadDefaults()
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		response := &types.ColorMappingResponse{
			Colors: make(map[string]*types.ColorResponse),
		}

		for nodename, color := range cm.GetColors() {
			response.Colors[nodename] = &types.ColorResponse{
				Red:   color.R,
				Green: color.G,
				Blue:  color.B,
				Alpha: color.A,
			}
		}

		cachedColormapping, err = json.Marshal(response)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(cachedColormapping)
}
