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

	crudex.Setup().WithUI(false)
	crudex.New[Car](db).OnRouter(app.Group("cars"))
	crudex.New[Driver](db).OnRouter(app.Group("drivers"))

	app.Run(":8080")
}
