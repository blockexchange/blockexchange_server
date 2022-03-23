package web

import (
	"blockexchange/types"
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

	w.Header().Set("Cache-Control", "max-age=345600")
	SendJson(w, response)
}
