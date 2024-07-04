package main

import (
	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := gin.New()
	db, _ := gorm.Open(sqlite.Open("sample.db"), &gorm.Config{})
	db.AutoMigrate(new(Car), new(Driver))

	crudex.
		Setup(app, db).           // setup crudex
		WithAutoScaffold(true).   // to generate the templates
		Add(crudex.New[Driver](), // create controllers and register them
						crudex.New[Car]()).
		Index("gen/index.html").    //and create index page
		OpenAPI("gen/openapi.json") // and create openapi spec

	app.HTMLRender = crudex.NewRenderer() // attach the crudex renderer
	app.Run(":8080")
}
