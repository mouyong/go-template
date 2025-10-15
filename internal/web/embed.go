package web

import "embed"

//go:embed build
var BuildFS embed.FS

//go:embed build/index.html
var IndexPage []byte
