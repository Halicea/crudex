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

    // the scaffold map contains all the scaffold templates and is used to generate the model templates
	ScaffoldMap() IScaffoldMap

	// Which template directories to scan for templates
	TemplateDirs() []string

	// the layout to use on the templates for full page rendering
	LayoutName() string

	// a function that is used to supply the layout with data
	LayoutDataFunc() func(c *gin.Context, data gin.H)

	// if true the layout will be used if the request is not an Htmx request, otherwise the template will be rendered without the layout
	EnableLayoutOnNonHxRequest() bool


    // if true the Scaffold templates will be exported to the 'scaffolds' directory
	ExportScaffolds() bool
}

type IRouter interface {
	gin.IRoutes
	Group(string, ...gin.HandlerFunc) *gin.RouterGroup
	BasePath() string
}

// IScaffoldMap is an interface that defines the scaffold map that is used to generate the model templates
type IScaffoldMap interface {
    // All returns a map of scaffolded template functions
    
    // These functions return the string representation of the scaffolded template (usually loaded from a file)
    // They are used to scaffold the templates for any given model
	All() map[string]func() string

    // Get returns a scaffolded template function by name
	Get(name string) func() string

    // Returns the function map that is passed to the template engine when generating the templates from the scaffolded templates
	FuncMap() template.FuncMap
}
