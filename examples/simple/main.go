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
	db.AutoMigrate(&Car{}, &Driver{})

	crudex.
		Setup(app, db).         // setup crudex
		WithAutoScaffold(true). // to generate the templates
		Add(                    // create controllers and register them
			crudex.New[Driver](),
			crudex.New[Car]()).
		Index("gen/index.html") //and create index page

	app.HTMLRender = crudex.NewRenderer() // attach the crudex renderer
	app.Run(":8080")
}
