package crudex

import (
	"fmt"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func NewRenderer() multitemplate.Renderer {
	return loadTemplates(config.TemplateDirs()...)
}

// loadTemplates is a helper function that loads the templates from the given directories
func loadTemplates(templatesDirs ...string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	if gin.IsDebugging() {
		fmt.Fprint(gin.DefaultWriter, "Loading templates from: ", templatesDirs, "\n")
	}
	for _, templatesDir := range templatesDirs {
		files, err := filepath.Glob(filepath.Join(templatesDir, "*.html"))
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			name := filepath.Base(file)
			if gin.IsDebugging() {
				fmt.Fprint(gin.DefaultWriter, "Loading template: ", file, " with name ", name, "\n")
			}
			r.AddFromFiles(name, file)
		}
	}
	return r
}
