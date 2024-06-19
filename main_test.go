package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/42dotmk/crudex/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
func TestHomepageHandler(t *testing.T) {

	app := setupApp()
	db := database.Connect()
	setupRenderer(app)
	registerRoutes(app, db)

	req, _ := http.NewRequest("GET", "/tags/", nil)
    req.Header.Add("Accept", "application/json")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	mockResponse := "hello world"
	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

