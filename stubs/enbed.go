package stubs

import "embed"

// just embed all swagger-ui static file to io.FS

//go:embed dist/*
var Static embed.FS

//go:embed favicon.ico
var Favicon []byte

//go:embed google.png microsoft.png
var Image embed.FS

//go:embed list.html detail.html
var Template embed.FS

//go:embed conf.toml.example
var ConfExample string
