package crudex

import (
	"html/template"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	ginrender "github.com/gin-gonic/gin/render"
)


func TestRender_Html(t *testing.T) {
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
	render(c, gin.H{"test": "test"}, "test")

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	if w.Body.String() != exp {
		t.Errorf("Expected %s, got %s", exp, w.Body.String())
	}
}

func TestRender_Json(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	c.Request = req
	c.Request.Header.Set("Accept", "application/json")
	render(c, gin.H{"test": "test"}, "test")
	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	if w.Body.String() != "{\"test\":\"test\"}" {
		t.Errorf("Expected {\"test\":\"test\"}, got %s", w.Body.String())
	}
}
