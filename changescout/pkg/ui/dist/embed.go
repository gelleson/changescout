package dist

import "embed"

var (
	//go:embed all:dist
	DistFS embed.FS
)
