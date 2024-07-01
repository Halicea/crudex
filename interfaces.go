package crudex

import (
	"text/template"

	"github.com/gin-gonic/gin"
)

// ICrudCtrl is an interface that defines the basic CRUD operations for a model
type ICrudCtrl interface {
	BasePath() string
	GetModelName() string

	List(c *gin.Context)
	Details(c *gin.Context)
	Form(c *gin.Context)
	Upsert(c *gin.Context)
	Delete(c *gin.Context)
}

// IConfig is an interface that defines the configuration for the crudex package
type IConfig interface {

	// how to create the templates
	ScaffoldStrategy() ScaffoldStrategy

	// where to place the scaffolded templates
	ScaffoldRootDir() string

	ScaffoldMap() IScaffoldMap

	// Which template directories to scan for templates
	TemplateDirs() []string

	// the layout to use on the templates for full page rendering
	LayoutName() string

	// a function that is used to supply the layout with data
	LayoutDataFunc() func(c *gin.Context, data gin.H)

	// if true the layout will be used if the request is not an Htmx request, otherwise the template will be rendered without the layout
	EnableLayoutOnNonHxRequest() bool

	ExportScaffolds() bool
}

type IRouter interface {
	gin.IRoutes
	Group(string, ...gin.HandlerFunc) *gin.RouterGroup
	BasePath() string
}

// IScaffoldMap is an interface that defines the scaffold map that is used to generate the model templates
type IScaffoldMap interface {
	All() map[string]func() string
	Get(name string) func() string
	FuncMap() template.FuncMap
}
