package crudex

import "github.com/gin-gonic/gin"

// NewConfig returns a new renderer Config with the default values
func NewConfig() *Config {
	return &Config{
		DefaultLayout:              "index.html",
		LayoutDataFunc:             nil,
		EnableLayoutOnNonHxRequest: true,
		TemplatesDirs:              []string{"templates"},
	}
}

// Config is used to configure the renderer and set the default layout and layout data function to use
// This config is than used by the Render function to render the data with the correct layout and data
type Config struct {
	//the layout to use on the templates for full page rendering
	DefaultLayout string

	//a function that is used to supply the layout with data
	LayoutDataFunc func(data gin.H)

	//if true the layout will be used even if the request is not an Htmx request, otherwise the template will be rendered without the layout
	EnableLayoutOnNonHxRequest bool

	// Which template directories to scan for templates
	TemplatesDirs []string

	// Templates suffix
	TemplateSuffix string
}

// the layout to use on the templates for full page rendering
func (c *Config) WithLayout(layout string) *Config {
	c.DefaultLayout = layout
	return c
}

// a function that is used to supply the layout with data
func (c *Config) WithLayoutDataFunc(layoutDataFunc func(data gin.H)) *Config {
	c.LayoutDataFunc = layoutDataFunc
	return c
}

// WithTemplateDirs sets the template directories to scan for templates when setting up the renderer
func (c *Config) WithTemplateDirs(dirs ...string) *Config {
	c.TemplatesDirs = dirs
	return c
}

// if true the layout will be used even if the request is not an Htmx request, otherwise the template will be rendered without the layout
func (c *Config) WithEnableLayoutOnNonHxRequest(enableLayoutOnNonHxRequest bool) *Config {
	c.EnableLayoutOnNonHxRequest = enableLayoutOnNonHxRequest
	return c
}

