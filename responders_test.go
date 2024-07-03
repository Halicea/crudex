package crudex

import (
	"html/template"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	ginrender "github.com/gin-gonic/gin/render"
)

func TestRespond_Html(t *testing.T) {
	tmpl, err := template.New("test").Parse(`<label>{{.test}}</label>`)

	if err != nil {
		t.Errorf("Error parsing template: %s", err)
	}
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, e := gin.CreateTestContext(w)
	e.HTMLRender = ginrender.HTMLProduction{
		Template: tmpl,
		Delims: ginrender.Delims{
			Left:  "{{",
			Right: "}}",
		},
	}
	req := httptest.NewRequest("GET", "/", nil)
	c.Request = req

	c.Request.Header.Set("Accept", "text/html")

	const exp = `<label>test</label>`
	_respond(c, gin.H{"test": "test"}, "test", "", &ResponseCapabilities{
		UI: true,
	})

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	if w.Body.String() != exp {
		t.Errorf("Expected %s, got %s", exp, w.Body.String())
	}
}

func TestRespond_Json(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	c.Request = req
	c.Request.Header.Set("Accept", "application/json")
	_respond(c, gin.H{"test": "test"}, "test", "", &ResponseCapabilities{
		UI:  true,
		API: true,
	})

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	if w.Body.String() != "{\"test\":\"test\"}" {
		t.Errorf("Expected {\"test\":\"test\"}, got %s", w.Body.String())
	}
}

func TestRespond_JsonIfOnlyThatCapabilityIsEnabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	c.Request = req
	c.Request.Header.Set("Accept", "test/html")
	_respond(c, gin.H{"test": "test"}, "test", "", &ResponseCapabilities{
		UI:  false,
		API: true,
	})
	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	if w.Body.String() != "{\"test\":\"test\"}" {
		t.Errorf("Expected {\"test\":\"test\"}, got %s", w.Body.String())
	}
}
func TestRender_ErrorIfNoUICapabilityIsMatched(t *testing.T) {
	c, w := faker()
	c.Request.Header.Set("Accept", "text/html")
	capabilities := ResponseCapabilities{API: false, UI: false}
    expected := "Error: No capability to respond for header Accept: text/html"
	_respond(c, nil, "test", "", &capabilities)

	if w.Code != 400 {
		t.Errorf("Expected 400, got %d", w.Code)
	}
	if w.Body.String() != expected {
		t.Errorf("Expected \"%s\", got %s", expected, w.Body.String())
	}

}
func TestRespond_ErrorIfNoAPICapabilityIsEnabled(t *testing.T) {
	c, w := faker()
	c.Request.Header.Set("Accept", "application/json")
    expected := "Error: No capability to respond for header Accept: application/json"
	assert := func() {
		if w.Code != 400 {
			t.Errorf("Expected 400, got %d", w.Code)
		}
		if w.Body.String() != expected {
            t.Errorf("Expected \"%s\", got %s", expected, w.Body.String())
		}
	}

	capabilities := ResponseCapabilities{API: false, UI: false}
	_respond(c, nil, "test", "", &capabilities)
	assert()
}

func TestRespond_ErrorIfInvalidAcceptHeader(t *testing.T) {
	c, w := faker()
	c.Request.Header.Set("Accept", "application/xml")
	expected := "Error: No capability to respond for header Accept: application/xml"
	_respond(c, nil, "test", "", &ResponseCapabilities{API: true, UI: true})
	if w.Code != 400 {
		t.Errorf("Expected 400, got %d", w.Code)
	}
	if w.Body.String() != expected {
		t.Errorf("Expected \"%s\", got %s", expected, w.Body.String())
	}
}


func TestRespond_ShouldOKIfInvalidAcceptHeaderAndOnlyOneCapabilityIsEnabled(t *testing.T) {
	c, w := faker()
	c.Request.Header.Set("Accept", "invalidheader")
    expected := "{\"test\":\"test\"}"
	_respond(c, gin.H{"test": "test"}, "test", "", &ResponseCapabilities{API: true, UI: false})
	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	if w.Body.String() != expected {
		t.Errorf("Expected \"%s\", got %s", expected, w.Body.String())
	}
}

func faker() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	c.Request = req
	return c, w
}
