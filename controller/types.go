package controller

import (
	"blockexchange/types"
)

type RenderData struct {
	Claims             *types.Claims
	BaseURL            string
	AdditionalMetaTags map[string]string
	Data               any
	IsAdmin            bool
}
