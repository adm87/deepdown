package data

import "embed"

//go:embed embedded/*
var EmbeddedFS embed.FS
