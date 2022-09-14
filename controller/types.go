package controller

import (
	"blockexchange/core"
	"blockexchange/types"
)

type RenderData struct {
	Claims             *types.Claims
	BaseURL            string
	Config             *core.Config
	AdditionalMetaTags map[string]string
	Data               any
	IsAdmin            bool
}
