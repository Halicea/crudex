package main

import (
	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex"
)

func main() {
	app := gin.New() // create a new gin app
	db := InitDb()   // open gorm db

	crudex.Setup(app, db).
		WithAutoScaffold(true).
        WithScaffoldStrategy(crudex.ScaffoldStrategyAlways).
		Add(crudex.New[Driver](), crudex.New[Car]()).
		Index("gen/index.html") //and create index page

	app.HTMLRender = crudex.NewRenderer() // attach the crudex renderer


    if err:=app.Run(":8080"); err != nil {
        panic(err)
    }
}
