package stubs

import "embed"

// just embed all swagger-ui static file to io.FS

//go:embed dist/*
var Static embed.FS

//go:embed favicon.ico
var Favicon []byte
