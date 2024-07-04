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
	db.AutoMigrate(&Car{}, &Driver{})

	c.Setup(app, db).WithAutoScaffold(true).
        WithScaffoldStrategy(c.ScaffoldStrategyIfNotExists).
		Add(
            c.New[Driver](), 
            c.New[Car]()).
		Index("gen/index.html")

	app.HTMLRender = c.NewRenderer()
	app.Run(":8080")
}

type Car struct {
	c.BaseModel
	Name        string `crud-input:"text" crud-placeholder:"Enter name"`
	License     string `crud-input:"html" crud-placeholder:"Enter the license plate"`
	Description string `crud-input:"wysiwyg" crud-placeholder:"Describe it"`
	Year        int    `crud-input:"number" crud-placeholder:"Model year of the car"`
}

type Driver struct {
	c.BaseModel
	Name  string
	CarID uint
	Car   Car `gorm:"foreignKey:CarID"`
}
