package crudex

import (
	"text/template"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	// HasUI returns true if the response should be rendered as a UI
	HasUI() bool

	// HasAPI returns true if the response should be rendered as an API
	HasAPI() bool

    // DefaultDb returns the default database connection
    DefaultDb() *gorm.DB

    // DefaultRouter returns the default router
    DefaultRouter() IRouter


    // AutoScaffold returns if every controller will scaffold it's ui automatically
    AutoScaffold() bool
}

// IResponseCapabilities is an interface that defines the capabilities of the response
//
// IConfig is already compliant with this interface
type IResponseCapabilities interface {
	// HasUI returns true if the response should be rendered as a UI
	HasUI() bool
	// HasAPI returns true if the response should be rendered as an API
	HasAPI() bool
	// EnableLayoutOnNonHxRequest returns true if the layout should be used even if the request is not an Htmx request
	EnableLayoutOnNonHxRequest() bool
}

// IRouter is an interface that defines the router that is used to scaffold the model templates
//
// It extends the gin.IRoutes interface but adds the BasePath method that returns the base path of the router.
//
// Note: the `RouterGroup` struct has a `BasePath()` method already, so there is no need to implement it
type IRouter interface {
	gin.IRoutes
	Group(string, ...gin.HandlerFunc) *gin.RouterGroup
	BasePath() string
}

// IScaffoldMap is an interface that defines the scaffold map that is used to generate the model templates
type IScaffoldMap interface {
	// All returns a map of scaffolded template functions
	//
	// These functions return the string representation of the scaffolded template (usually loaded from a file)
	// They are used to scaffold the templates for any given model
	All() map[string]func() string

	// Get returns a scaffolded template function by name
	Get(name string) func() string

	// Returns the function map that is passed to the template engine when generating the templates from the scaffolded templates
	FuncMap() template.FuncMap

	// Export exports the scaffolded templates to the file system in the `scaffolds` directory.
	//
	// If the `forceIfExists` parameter is true, it will overwrite any existing scaffold templates
	Export(forceIfExists bool) error
}
