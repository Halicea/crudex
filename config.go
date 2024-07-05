package crudex

import (
	"flag"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex/scaffolds"
	_ "github.com/pboyd04/godata/middleware"
	"gorm.io/gorm"

    _ "github.com/pboyd04/godata/filter/parser/gorm"
	odata "github.com/pboyd04/godata/middleware"
)

type Config struct {
	scaffoldCreateStrategy ScaffoldStrategy
	// where to place the scaffolded templates
	scaffoldRootDir string
	scaffoldMap     IScaffoldMap

	// Which template directories to scan for templates
	templateDirs []string

	//the layout to use on the templates for full page rendering
	layoutName string

	//if true the layout will be used even if the request is not an Htmx request,
	//otherwise the template will be rendered without the layout
	enableLayoutOnNonHxRequest bool

	//a function that is used to supply the layout with data
	layoutDataFunc func(c *gin.Context, data gin.H)

	//Enable API requests
	apiEnabled bool

	//Enable UI requests
	uiEnabled bool

	//Default database
	defaultDb *gorm.DB

	//Default router
	defaultRouter IRouter

	// list of the controller registered with the configuration
	controllers *ControllerList

	// wether to auto scaffold the templates when a new controller is created
	autoScaffold bool
}

// NewConfig creates a new configuration crud configuration containing all the defaults
func Setup(router IRouter, db *gorm.DB) *Config {
	router.Use(odata.NewOdataMiddleware(nil).GinMiddleware)
	return NewConfig().
		WithDefaultDb(db).
		WithDefaultRouter(router).
		SetAsDefault()
}

// config is the default configuration
var config IConfig = NewConfig()

// GetConfig returns the default crudex configuration for the package
func GetConfig() IConfig {
	return config
}

func NewConfig() *Config {
	return &Config{
		scaffoldCreateStrategy: ScaffoldStrategyIfNotExists,
		scaffoldRootDir:        "gen",
		scaffoldMap:            scaffolds.New(),
		autoScaffold:           true,

		templateDirs:               []string{"gen", "templates"},
		layoutName:                 "index.html",
		enableLayoutOnNonHxRequest: true,
		layoutDataFunc:             nil,

		apiEnabled: true,
		uiEnabled:  true,

		defaultDb:     nil,
		defaultRouter: nil,

		controllers: &ControllerList{},
	}
}

func (conf *Config) String() string {
	return fmt.Sprintf(`
#############################################
Crudex Configuration:
    Scaffold Strategy: %s
    Scaffold Root Dir: %s
    Scaffold Map: %s
    Template Dirs: %s
    Layout Name: %s
    Enable Layout On Non Hx Requests: %t
    API Enabled: %t
    UI Enabled: %t
#############################################`,
		conf.scaffoldCreateStrategy,
		conf.scaffoldRootDir,
		conf.scaffoldMap,
		conf.templateDirs,
		conf.layoutName,
		conf.enableLayoutOnNonHxRequest,
		conf.apiEnabled,
		conf.uiEnabled)
}

// how to create the templates
func (conf *Config) ScaffoldStrategy() ScaffoldStrategy {
	return conf.scaffoldCreateStrategy
}

func (conf *Config) ScaffoldMap() IScaffoldMap {
	return conf.scaffoldMap
}

// where to create the templates
func (conf *Config) ScaffoldRootDir() string {
	return conf.scaffoldRootDir
}

// Which template directories to scan for templates
func (conf *Config) TemplateDirs() []string {
	return conf.templateDirs
}

// the layout to use on the templates for full page rendering
func (conf *Config) LayoutName() string {
	return conf.layoutName
}

// if true the layout will be used even if the request is not an Htmx request, otherwise the template will be rendered without the layout
func (conf *Config) EnableLayoutOnNonHxRequest() bool {
	return conf.enableLayoutOnNonHxRequest
}

// a function that is used to supply the layout with data
func (conf *Config) LayoutDataFunc() func(c *gin.Context, data gin.H) {
	return conf.layoutDataFunc
}

// HasUI returns true if the configuration has the UI enabled
func (conf *Config) HasUI() bool {
	return conf.uiEnabled
}

// HasAPI returns true if the configuration has the API enabled
func (conf *Config) HasAPI() bool {
	return conf.apiEnabled
}

// DefaultDb returns the default database connection
func (conf *Config) DefaultDb() *gorm.DB {
	return conf.defaultDb
}

// DefaultRouter returns the default router
func (conf *Config) DefaultRouter() IRouter {
	return conf.defaultRouter
}

// Index creates a simple index page that lists all the controllers registered with the configuration
func (conf *Config) Index(template string) *Config {
	conf.Controllers().
		Index(conf.DefaultRouter(), template, conf)
	return conf
}

// Index creates OpenAPI spec
func (conf *Config) OpenAPI(template string) *Config {
	conf.Controllers().
		OpenAPI(conf.DefaultRouter(), template, conf)
	return conf
}

// Controllers returns the list of controllers registered with the configuration
func (conf *Config) Controllers() *ControllerList {
	return conf.controllers
}

// Add adds the controllers to the configuration
func (conf *Config) Add(ctrls ...ICrudCtrl) *Config {
	conf.controllers.Add(ctrls...)
	return conf
}

func (conf *Config) AutoScaffold() bool {
	return conf.autoScaffold
}

// WithScaffoldStrategy sets the strategy to use when creating the scaffolded templates
// The default is ScaffoldCreateAlways, options are ScaffoldCreateAlways, ScaffoldCreateIfNotExist, ScaffoldCreateNever
// This option is not used at the moment
func (conf *Config) WithScaffoldStrategy(value ScaffoldStrategy) *Config {
	conf.scaffoldCreateStrategy = value
	return conf
}

// WithScaffoldRootDir sets the root directory where the scaffolded templates will be placed
func (conf *Config) WithScaffoldRootDir(value string) *Config {
	conf.scaffoldRootDir = value
	return conf
}

// the layout to use on the templates for full page rendering
func (c *Config) WithLayoutName(layoutName string) *Config {
	c.layoutName = layoutName
	return c
}

// WithLayoutDataFunc is function that is used to supply the layout with the data needed to render the layout
func (c *Config) WithLayoutDataFunc(layoutDataFunc func(c *gin.Context, data gin.H)) *Config {
	c.layoutDataFunc = layoutDataFunc
	return c
}

// WithTemplateDirs sets the template directories to scan for templates when setting up the renderer
func (c *Config) WithTemplateDirs(dirs ...string) *Config {
	c.templateDirs = dirs
	return c
}

// if true the layout will be used even if the request is not an Htmx request, otherwise the template will be rendered without the layout
func (c *Config) WithEnableLayoutOnNonHxRequest(enableLayoutOnNonHxRequest bool) *Config {
	c.enableLayoutOnNonHxRequest = enableLayoutOnNonHxRequest
	return c
}

// SetAsDefault sets the current configuration as the default configuration
func (conf *Config) SetAsDefault() *Config {
	config = conf
	return conf
}

// WithScaffoldMap sets the scaffold map that will be used to generate the scaffolded templates
func (conf *Config) WithScaffoldMap(scaffoldMap IScaffoldMap) *Config {
	conf.scaffoldMap = scaffoldMap
	return conf
}

// WithAPI sets the configuration to enable the API endpoints
func (conf *Config) WithAPI(value bool) *Config {
	conf.apiEnabled = value
	return conf
}

// WithUI sets the configuration to enable the UI endpoints
func (conf *Config) WithUI(value bool) *Config {
	conf.uiEnabled = value
	return conf
}

// WithDefaultRouter sets the default router to use when creating the controllers
func (conf *Config) WithDefaultRouter(router IRouter) *Config {
	conf.defaultRouter = router
	return conf
}

// WithDefaultDb sets the default database to use when creating the controllers
func (conf *Config) WithDefaultDb(db *gorm.DB) *Config {
	conf.defaultDb = db
	return conf
}

// WithAutoScaffold sets the configuration to auto scaffold the templates when a new controller is created
func (conf *Config) WithAutoScaffold(value bool) *Config {
	conf.autoScaffold = value
	return conf
}

// WithCommandLineArgs sets the configuration from the command line arguments
func (conf *Config) WithCommandLineArgs(args []string) *Config {
	var templateDirs string
	var layout string
	var hxAware string
	var scaffoldExportBases string
	var scaffoldDir string
	var scaffoldStrategy string
	flags := flag.NewFlagSet("crudex", flag.PanicOnError)
	flags.StringVar(&templateDirs, "crud-template-dirs", "", "Template directories")
	flags.StringVar(&layout, "crud-layout", "", "The main layout to use for the hxAware rendering")
	flags.StringVar(&hxAware, "crud-hx-aware", "", "Template directories")
	flags.StringVar(&scaffoldExportBases, "crud-export-bases", "", "Wether to export the templates needed to generate the scaffolds")
	flags.StringVar(&scaffoldDir, "crud-scaffold-dir", "", "Where to export the generated templates")
	flags.StringVar(&scaffoldStrategy, "crud-strategy", "", `When to export the templates
    - always: Exports even if the files already exist(any changes will be overwritten)
    - newonly: Exports a template only if the template file does not already exist
    - never: No templates will be exported`)

	if error := flags.Parse(args); error != nil {
		panic(error)
	}
	if templateDirs != "" {
		conf.WithTemplateDirs(strings.Split(templateDirs, ",")...)
	}
	if layout != "" {
		conf.WithLayoutName(layout)
	}
	if hxAware != "" {
		conf.WithEnableLayoutOnNonHxRequest(hxAware == "true")
	}
	if scaffoldDir != "" {
		conf.WithScaffoldRootDir(scaffoldDir)
	}
	if scaffoldStrategy != "" {
		switch scaffoldStrategy {
		case CmdArgStrategyAlways:
			conf.WithScaffoldStrategy(ScaffoldStrategyAlways)
		case CmdArgStrategyIfNotExists:
			conf.WithScaffoldStrategy(ScaffoldStrategyIfNotExists)
		case CmdArgStrategyNever:
			conf.WithScaffoldStrategy(ScaffoldStrategyNever)
		default:
			panic("Invalid strategy")
		}
	}
	return conf
}
