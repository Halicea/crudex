package crudex

import (
	"flag"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex/scaffolds"
)

type Config struct {
	scaffoldCreateStrategy ScaffoldStrategy
	// where to place the scaffolded templates
	scaffoldRootDir string
	exportScaffolds bool
	scaffoldMap     IScaffoldMap

	// Which template directories to scan for templates
	templateDirs []string

	//the layout to use on the templates for full page rendering
	layoutName string

	//if true the layout will be used even if the request is not an Htmx request, otherwise the template will be rendered without the layout
	enableLayoutOnNonHxRequest bool

	//a function that is used to supply the layout with data
	layoutDataFunc func(c *gin.Context, data gin.H)
}

// NewConfig creates a new configuration crud configuration containing all the defaults
func Setup() *Config {
	return NewConfig().SetAsDefault()
}

func NewConfig() *Config {
	return &Config{
		scaffoldCreateStrategy:      SCAFFOLD_ALWAYS,
		scaffoldRootDir:            "gen",
		scaffoldMap:                scaffolds.New(),
		exportScaffolds:            false,
		layoutName:                 "index.html",
		enableLayoutOnNonHxRequest: true,
		layoutDataFunc:             nil,
		templateDirs:               []string{"gen", "templates"},
	}
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

// ExportScaffolds gets if the scaffold templates should be exported to the file system
func (self *Config) ExportScaffolds() bool {
	return self.exportScaffolds
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

// WithExportScaffolds gets if the scaffold templates should be exported to the file system
func (self *Config) WithExportScaffolds(export bool) *Config {
	self.exportScaffolds = export
	return self
}


// SetAsDefault sets the current configuration as the default configuration
func (self *Config) SetAsDefault() (config *Config) {
	config = self
	return self
}

// WithScaffoldMap sets the scaffold map that will be used to generate the scaffolded templates
func (self *Config) WithScaffoldMap(scaffoldMap IScaffoldMap) *Config {
    self.scaffoldMap = scaffoldMap
    return self
}

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
	if scaffoldExportBases != "" {
		self.WithExportScaffolds(scaffoldExportBases == "true")
	}
	if scaffoldDir != "" {
		self.WithScaffoldRootDir(scaffoldDir)
	}
	if scaffoldStrategy != "" {
		switch scaffoldStrategy {
		case CmdArgStrategyAlways:
			self.WithScaffoldStrategy(SCAFFOLD_ALWAYS)
		case CmdArgStrategyIfNotExists:
			self.WithScaffoldStrategy(SCAFFOLD_IF_NOT_EXISTS)
		case CmdArgStrategyNever:
			self.WithScaffoldStrategy(SCAFFOLD_NEVER)
		default:
			panic("Invalid strategy")
		}
	}
	return self
}
