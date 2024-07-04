package crudex

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ControllerList []ICrudCtrl

// ScaffoldIndex creates a simple index page that lists all the controllers
func (list *ControllerList) Index(r IRouter, templateFile string, conf IConfig) *ControllerList {
	arr := []ICrudCtrl(*list)
	GenLayout(templateFile, arr)
	r.GET("/", func(c *gin.Context) {
		data := gin.H{"Path": r.BasePath()}
		template := filepath.Base(templateFile)
		RespondWithConfig(c, data, template, conf)
	})
	return list
}

func (list *ControllerList) OpenAPI(r IRouter, templateFile string, conf IConfig) *ControllerList {
	arr := []ICrudCtrl(*list)
	GenOpenAPI(templateFile, arr)
	r.GET("/openapi", func(c *gin.Context) {
		c.File(templateFile)
	})
	return list
}

func (list *ControllerList) Add(r ...ICrudCtrl) *ControllerList {
	*list = append(*list, r...)
	return list
}
