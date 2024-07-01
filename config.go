package crudex

import (
	"flag"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex/scaffolds"
)
//go:generate stringer -type=ScaffoldStrategy
type ScaffoldStrategy int
const (
	// SCAFFOLD_ALWAYS will always scaffold the model templates
	SCAFFOLD_ALWAYS ScaffoldStrategy = iota
	// SCAFFOLD_IF_NOT_EXISTS will only scaffold the model templates if they do not exist
	SCAFFOLD_IF_NOT_EXISTS
	// SCAFFOLD_NEVER will never scaffold the model templates
	SCAFFOLD_NEVER
)


type IConfig interface {

	// how to create the templates
	ScaffoldStrategy() ScaffoldStrategy

	// where to place the scaffolded templates
	ScaffoldRootDir() string

	// which functions to regsiter for the screation of scaffoled templates
	ScaffoldFuncs() template.FuncMap

	// ListScaffold defines the scaffold template that generages the list[T] template
	//
	// Note: the delimiters are `[[` `]]`
	ListScaffold() string
	// DetailScaffold defines the scaffold template that generages the defail[T] template
	//
	// Note: the delimiters are `[[` `]]`
	DetailScaffold() string
	// FormScaffold defines the scaffold template that generages the form[T] template
	//
	// Note: the delimiters are `[[` `]]`
	FormScaffold() string

	// LayoutScaffold defines the scaffold template that generates the layout template used for listing all the models
	//
	// Note: the delimiters are `[[` `]]`
	LayoutScaffold() string

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

type Config struct {
	scafoldCreateStrategy ScaffoldStrategy
	// where to place the scaffolded templates
	scaffoldRootDir string
	exportScaffolds bool
	// The func map passed to the templates so they can use the functions defined
	scaffoldFuncs  template.FuncMap
	listScaffold   func() string
	detailScaffold func() string
	formScaffold   func() string
	layoutScaffold func() string

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
func DefaultConfig() *Config {
    return NewConfig().SetAsDefault()
}
func NewConfig() *Config {
	return &Config{
		scafoldCreateStrategy: SCAFFOLD_ALWAYS,
		scaffoldRootDir:       "gen",
		scaffoldFuncs: template.FuncMap{
			"RenderTypeInput": RenderInputType,
		},
		exportScaffolds: true,

		layoutScaffold: func() string { return readContentsOrDefault("scaffolds/layout.html", scaffolds.Layout, true) },
		listScaffold:   func() string { return readContentsOrDefault("scaffolds/list.html", scaffolds.List, true) },
		detailScaffold: func() string { return readContentsOrDefault("scaffolds/detail.html", scaffolds.Detail, true) },
		formScaffold:   func() string { return readContentsOrDefault("scaffolds/form.html", scaffolds.Form, true) },

		layoutName:                 "index.html",
		enableLayoutOnNonHxRequest: true,
		layoutDataFunc:             nil,
		templateDirs:               []string{"gen", "templates"},
	}
}

// how to create the templates
func (self *Config) ScaffoldStrategy() ScaffoldStrategy {
	return self.scafoldCreateStrategy
}

// where to create the templates
func (self *Config) ScaffoldRootDir() string {
	return self.scaffoldRootDir
}

// based on what template to create the list[T] template
//
// Note: the delimiters are `[[` `]]`
func (self *Config) ListScaffold() string {
	return self.listScaffold()
}

// based on what template to create the details[T] template
//
// Note: the delimiters are `[[` `]]`
func (self *Config) DetailScaffold() string {
	return self.detailScaffold()
}

// based on what template to create the form[T] template used for edit and create
//
// Note: the delimiters are `[[` `]]`
func (self *Config) FormScaffold() string {
	return self.formScaffold()
}

// based on what template to create the layout template used for listing all the models
//
// Note: the delimiters are `[[` `]]`
func (self *Config) LayoutScaffold() string {
	return self.layoutScaffold()
}

// The func map passed to the templates so they can use the functions defined
func (self *Config) ScaffoldFuncs() template.FuncMap {
	return self.scaffoldFuncs
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
	self.scafoldCreateStrategy = value
	return self
}

// WithScaffoldRootDir sets the root directory where the scaffolded templates will be placed
func (self *Config) WithScaffoldRootDir(value string) *Config {
	self.scaffoldRootDir = value
	return self
}

// WithListScaffold sets the scaffold template that generages the list[T] template
func (self *Config) WithListScaffold(value func() string) *Config {
	self.listScaffold = value
	return self
}

// WithDetailScaffold sets the scaffold template function that generages the detail[T] template
func (self *Config) WithDetailScaffold(value func() string) *Config {
	self.detailScaffold = value
	return self
}

// WithFormScaffold sets the scaffold template function that generages the form[T] template
func (self *Config) WithFormScaffold(value func() string) *Config {
	self.formScaffold = value
	return self
}

// WithLayoutScaffold sets the scaffold template function that generates the layout template used for listing all the models
func (self *Config) WithLayoutScaffold(value func() string) *Config {
	self.layoutScaffold = value
	return self
}

// WithScaffoldFuncs sets the functions to be used when scaffolding the templates
func (self *Config) WithScaffoldFuncs(value template.FuncMap) *Config {
	self.scaffoldFuncs = value
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

const (
    ArgStrategyAlways    = "always"
    ArgStrategyIfNotExists = "newonly"
    ArgStrategyNever     = "never"
)

func (self *Config) SetAsDefault() (config *Config) {
    config = self
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
		case ArgStrategyAlways:
			self.WithScaffoldStrategy(SCAFFOLD_ALWAYS)
		case ArgStrategyIfNotExists:
			self.WithScaffoldStrategy(SCAFFOLD_IF_NOT_EXISTS)
		case ArgStrategyNever:
			self.WithScaffoldStrategy(SCAFFOLD_NEVER)
        default:
            panic("Invalid strategy")
		}
	}
	return self
}
