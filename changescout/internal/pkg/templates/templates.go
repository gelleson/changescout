package templates

import (
	_ "embed"
)

var (
	//go:embed diff-default-message.tpl
	DiffDefaultMessage string
)
