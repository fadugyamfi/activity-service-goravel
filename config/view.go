package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("view", map[string]any{
		// View rendering engine: html (Go templates)
		"engine": "html",

		// Path to view files
		"path": "resources/views",
	})
}
