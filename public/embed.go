package public

import (
	"embed"
)

// built webapp (npm run build)
//
//go:embed all:dist
var Webapp embed.FS

const Prefix = "dist"
