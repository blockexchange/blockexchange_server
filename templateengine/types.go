package templateengine

import (
	"blockexchange/types"
)

type RenderData struct {
	Claims  *types.Claims
	BaseURL string
	Data    any
	IsAdmin bool
}
