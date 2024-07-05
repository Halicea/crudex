package main

import (
	"github.com/gin-gonic/gin"
	c "github.com/halicea/crudex"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := gin.New()
	db, _ := gorm.Open(sqlite.Open("sample.db"), &gorm.Config{})
	db.AutoMigrate(new(Driver))
    conf:= c.Setup(app, db). // setup
        WithUI(true).
        WithAPI(true).
        WithAutoScaffold(true). // to generate the templates
        WithScaffoldStrategy(c.ScaffoldStrategyIfNotExists).
        Add(
			c.New[Driver](),
			c.New[Passenger](),
		).
		Index("gen/index.html") //and create index page

    conf.ScaffoldMap().Export(false)


	app.HTMLRender = c.NewRenderer() // attach the c renderer
	app.Run(":8080")
}

type Passenger struct {
	c.BaseModel
	Name     string
	Location string
}

type Driver struct {
	c.BaseModel
	Name    string
	Surname string
}
