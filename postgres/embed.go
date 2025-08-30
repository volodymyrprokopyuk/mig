package postgres

import (
	"embed"
)

//go:embed *.apply.sql *.revert.sql
var FS embed.FS
