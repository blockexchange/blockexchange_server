package web

import (
	"blockexchange/types"
	"encoding/json"
	"net/http"
)

func (api *Api) GetColorMapping(w http.ResponseWriter, r *http.Request) {
	response := &types.ColorMappingResponse{
		Colors: make(map[string]*types.ColorResponse),
	}

	for nodename, color := range api.ColorMapping.GetColors() {
		response.Colors[nodename] = &types.ColorResponse{
			Red:   color.R,
			Green: color.G,
			Blue:  color.B,
			Alpha: color.A,
		}
	}

	colormapping, err := json.Marshal(response)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	SendJson(w, colormapping)
}
