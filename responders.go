package crudex

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Respond is a function that renders a template with the given data
// If the request is request accepts application/json it will the data as json
// If the request is an Htmx request it will render the template with the data
// If the request is not an Htmx request it will use the layout to render the data
//   - The layout should be aware of the data that is passed to it and conditionally render that template
//
// it uses the default hxConfig to render the template
// See `RenderWithConfig` for more control over the rendering
func Respond(c *gin.Context, data gin.H, templateName string) {
	RespondWithConfig(c, data, templateName, GetConfig())
}

// RespondWithConfig is a function that renders a template with the given data and the render configuration
// See Render for more information on the rendering
// See Config for more information on the configuration
func RespondWithConfig(c *gin.Context, data gin.H, templateName string, conf IConfig) {
	if conf.LayoutDataFunc() != nil {
		config.LayoutDataFunc()(c, data)
	}
	_respond(c, data, templateName, conf.LayoutName(), conf)
}

// _respond is a helper function that renders the data based on the request accept header and the Hx-Request header
//
// If the request accepts application/json it will render the data as json
// If the request accepts text/html it will render the data as html
// If the request has the Hx-Request header set to true it will render the template with the data on full page load
//
// capabilites is an interface that defines the capabilities of the response
//   - HasUI() bool
//   - HasApi() bool
//   - EnableLayoutOnNonHxRequest() bool
//
// If the request does not match the capabilities it will write an error to the response
func _respond(c *gin.Context, data gin.H, templateName string, layout string, capabilites IResponseCapabilities) {
	var isNone = c.Request.Header.Get("Accept") == ""                    // no accept header, we default to json
	var isStar = strings.Contains(c.Request.Header.Get("Accept"), "*/*") // we default to html
	var isApi = strings.Contains(c.Request.Header.Get("Accept"), "application/json")
	var isUi = strings.Contains(c.Request.Header.Get("Accept"), "text/html")
	var isHxRequest = c.Request.Header.Get("Hx-Request") == "true"
	hasUI := capabilites.HasUI()
	hasAPI := capabilites.HasAPI()
	useLayoutOnFullPageLoad := capabilites.EnableLayoutOnNonHxRequest()
	layoutEnabled := layout != "" && useLayoutOnFullPageLoad

	switch {

	case hasAPI && (isApi || isNone || !hasUI):
		c.JSON(http.StatusOK, data)
	case hasUI && (isUi || isStar || !hasAPI):
		if isHxRequest || !layoutEnabled {
			data["IsLayoutEnabled"] = false
			c.HTML(http.StatusOK, templateName, data)
		} else {
			data["IsLayoutEnabled"] = true
			c.HTML(http.StatusOK, templateName, data)
		}
	default:
		err := fmt.Errorf("No capability to respond for header Accept: %s", c.Request.Header.Get("Accept"))
		c.String(http.StatusBadRequest, "Error: %s", err.Error())
	}
}
