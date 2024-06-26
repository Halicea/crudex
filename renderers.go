package crudex

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

// config is the default configuration for the Render function
var config IConfig = NewConfig()

// sets the configuration for the Render function and returns the renderer
func NewRenderer(options IConfig) multitemplate.Renderer {
	config = options
	return loadTemplates(config.TemplateDirs()...)
}

// Render is a function that renders a template with the given data
// If the request is request accepts application/json it will return the data as json
// If the request is an Htmx request it will render the template with the data
// If the request is not an Htmx request it will use the layout to render the data
//   - The layout should be aware of the data that is passed to it and conditionally render that template
//
// it uses the default hxConfig to render the template
// See `RenderWithConfig` for more control over the rendering
func Render(c *gin.Context, data gin.H, templateName string) {
	hxAwareRender(c, data, templateName,
		config.LayoutName(),
		config.EnableLayoutOnNonHxRequest())
}

// RenderWithConfig is a function that renders a template with the given data and the render configuration
// See Render for more information on the rendering
// See Config for more information on the configuration
func RenderWithConfig(c *gin.Context, data gin.H, templateName string, conf IConfig) {
	if conf.LayoutDataFunc() != nil {
		config.LayoutDataFunc()(c, data)
	}
	hxAwareRender(c, data, templateName,
		config.LayoutName(),
		conf.EnableLayoutOnNonHxRequest())
}

// loadTemplates is a helper function that loads the templates from the given directories
func loadTemplates(templatesDirs ...string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	if gin.IsDebugging() {
		fmt.Fprint(gin.DefaultWriter, "Loading templates from: ", templatesDirs, "\n")
	}
	for _, templatesDir := range templatesDirs {
		files, err := filepath.Glob(filepath.Join(templatesDir, "*.html"))
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

// RenderTypeInput is a helper function that renders an input based on the type of the field.
//
// This function is part of the default FuncMap that is passed to the scaffold templates.
// It is used in the form template to render the input fields for the model
func RenderTypeInput(modelName string, field reflect.StructField) string {
	switch field.Type.Kind() {
	case reflect.String:
		return fmt.Sprintf(`<input type="text" name="%s">{{.%s.%s}}</input>`, field.Name, modelName, field.Name)
	case reflect.Int, reflect.Float64:
		return fmt.Sprintf(`<input type="number" name="%s">{{.%s.%s}}</input>`, field.Name, modelName, field.Name)
	case reflect.Bool:
		return fmt.Sprintf(`<input type="checkbox" name="%s">{{.%s.%s}}</input>`, field.Name, modelName, field.Name)
	}
	panic(fmt.Sprintf("unsupported type: %s for field %s", field.Type.Kind().String(), field.Name))
}

// hxAwareRender is a helper function that renders the data based on the request accept header and the Hx-Request header
func hxAwareRender(c *gin.Context, data gin.H, templateName string, layout string, enableLayoutOnNonHxRequest bool) {
	if c.Request.Header.Get("Hx-Request") == "true" || !enableLayoutOnNonHxRequest {
		c.HTML(http.StatusOK, templateName, data)
	} else {
		render(c, data, layout)
	}
}

// render is a helper function that renders the data based on the request accept header
func render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data)
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
