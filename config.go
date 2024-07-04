package crudex

import (
	"flag"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex/scaffolds"
	"gorm.io/gorm"
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

	controllers *ControllerList

	autoScaffold bool
}

// NewConfig creates a new configuration crud configuration containing all the defaults
func Setup(router IRouter, db *gorm.DB) *Config {
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
		scaffoldCreateStrategy:     ScaffoldStrategyIfNotExists,
		scaffoldRootDir:            "gen",
		scaffoldMap:                scaffolds.New(),
		autoScaffold:               true,

		templateDirs:               []string{"gen", "templates"},
		layoutName:                 "index.html",
		enableLayoutOnNonHxRequest: true,
		layoutDataFunc:             nil,

		apiEnabled: true,
		uiEnabled:  true,

		defaultDb: nil,
		defaultRouter: nil,

		controllers: &ControllerList{},
	}
}

func (self *Config) String() string {
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
		self.scaffoldCreateStrategy,
		self.scaffoldRootDir,
		self.scaffoldMap,
		self.templateDirs,
		self.layoutName,
		self.enableLayoutOnNonHxRequest,
		self.apiEnabled,
		self.uiEnabled)
}

// how to create the templates
func (self *Config) ScaffoldStrategy() ScaffoldStrategy {
	return self.scaffoldCreateStrategy
}

func (self *Config) ScaffoldMap() IScaffoldMap {
	return self.scaffoldMap
}

// where to create the templates
func (self *Config) ScaffoldRootDir() string {
	return self.scaffoldRootDir
}

// Which template directories to scan for templates
func (self *Config) TemplateDirs() []string {
	return self.templateDirs
}

// the layout to use on the templates for full page rendering
func (self *Config) LayoutName() string {
	return self.layoutName
}

// if true the layout will be used even if the request is not an Htmx request, otherwise the template will be rendered without the layout
func (self *Config) EnableLayoutOnNonHxRequest() bool {
	return self.enableLayoutOnNonHxRequest
}

// a function that is used to supply the layout with data
func (self *Config) LayoutDataFunc() func(c *gin.Context, data gin.H) {
	return self.layoutDataFunc
}

// HasUI returns true if the configuration has the UI enabled
func (self *Config) HasUI() bool {
	return self.uiEnabled
}

// HasAPI returns true if the configuration has the API enabled
func (self *Config) HasAPI() bool {
	return self.apiEnabled
}

// DefaultDb returns the default database connection
func (self *Config) DefaultDb() *gorm.DB {
	return self.defaultDb
}

// DefaultRouter returns the default router
func (self *Config) DefaultRouter() IRouter {
	return self.defaultRouter
}

// Index creates a simple index page that lists all the controllers registered with the configuration
func (conf *Config) Index(template string) *ControllerList {
	return conf.Controllers().
		Index(conf.DefaultRouter(), template, conf)
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
func (self *Config) WithScaffoldStrategy(value ScaffoldStrategy) *Config {
	self.scaffoldCreateStrategy = value
	return self
}

// WithScaffoldRootDir sets the root directory where the scaffolded templates will be placed
func (self *Config) WithScaffoldRootDir(value string) *Config {
	self.scaffoldRootDir = value
	return self
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
func (self *Config) SetAsDefault() *Config {
	config = self
	return self
}

// WithScaffoldMap sets the scaffold map that will be used to generate the scaffolded templates
func (self *Config) WithScaffoldMap(scaffoldMap IScaffoldMap) *Config {
	self.scaffoldMap = scaffoldMap
	return self
}

// WithAPI sets the configuration to enable the API endpoints
func (self *Config) WithAPI(value bool) *Config {
	self.apiEnabled = value
	return self
}

// WithUI sets the configuration to enable the UI endpoints
func (self *Config) WithUI(value bool) *Config {
	self.uiEnabled = value
	return self
}

// WithDefaultRouter sets the default router to use when creating the controllers
func (self *Config) WithDefaultRouter(router IRouter) *Config {
	self.defaultRouter = router
	return self
}

// WithDefaultDb sets the default database to use when creating the controllers
func (self *Config) WithDefaultDb(db *gorm.DB) *Config {
	self.defaultDb = db
	return self
}

func (conf *Config) WithAutoScaffold(value bool) *Config {
	conf.autoScaffold = value
	return conf
}

// WithCommandLineArgs sets the configuration from the command line arguments
func (self *Config) WithCommandLineArgs(args []string) *Config {
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
		self.WithTemplateDirs(strings.Split(templateDirs, ",")...)
	}
	if layout != "" {
		self.WithLayoutName(layout)
	}
	if hxAware != "" {
		self.WithEnableLayoutOnNonHxRequest(hxAware == "true")
	}
	if scaffoldDir != "" {
		self.WithScaffoldRootDir(scaffoldDir)
	}
	if scaffoldStrategy != "" {
		switch scaffoldStrategy {
		case CmdArgStrategyAlways:
			self.WithScaffoldStrategy(ScaffoldStrategyAlways)
		case CmdArgStrategyIfNotExists:
			self.WithScaffoldStrategy(ScaffoldStrategyIfNotExists)
		case CmdArgStrategyNever:
			self.WithScaffoldStrategy(ScaffoldStrategyNever)
		default:
			panic("Invalid strategy")
		}
	}
	return self
}
