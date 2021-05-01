package static

import (
	_ "embed"
)

var (
	//go:embed embedded/version.txt
	Version string
	//go:embed embedded/builddate.txt
	BuildDate string
)
