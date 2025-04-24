package golang

import "embed"

//go:embed templates/*
var templatesFS embed.FS
