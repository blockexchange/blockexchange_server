package wasm

import "embed"

// luanti wasm client build for "Edit Online", see README.md
// large files are stored gzip-compressed (luanti.wasm.gz), statigz serves them under the original name
//
//go:embed all:runtime
var Runtime embed.FS

const Prefix = "runtime"
