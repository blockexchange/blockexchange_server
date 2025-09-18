package public

import (
	"embed"
)

//go:embed js/* pics/* index.html
//go:embed node_modules/bootswatch/dist/darkly/bootstrap.min.css
//go:embed node_modules/vue/dist/vue.global.js
//go:embed node_modules/vue/dist/vue.global.prod.js
//go:embed node_modules/vue-router/dist/vue-router.global.js
//go:embed node_modules/vue-router/dist/vue-router.global.prod.js
//go:embed node_modules/@fortawesome/fontawesome-free/css/all.min.css
//go:embed node_modules/@fortawesome/fontawesome-free/webfonts/*
//go:embed node_modules/marked/lib/marked.umd.js
//go:embed node_modules/dompurify/dist/purify.min.js
//go:embed node_modules/three/build/three.module.min.js
//go:embed node_modules/three/build/three.core.min.js
//go:embed node_modules/fflate/esm/browser.js
var Webapp embed.FS
