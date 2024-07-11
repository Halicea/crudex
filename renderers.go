package crudex

import (
	"fmt"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func NewRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
    templateDirs := GetConfig().TemplateDirs()
	if gin.IsDebugging() {
		fmt.Fprint(gin.DefaultWriter, "Loading templates from: ", templateDirs, "\n")
	}

	for _, td := range templateDirs {
		files, err := filepath.Glob(filepath.Join(td, "*.html"))
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
