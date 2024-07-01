package main

import (
	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := gin.New()                                             //create gin app
	db, _ := gorm.Open(sqlite.Open("sample.db"), &gorm.Config{}) //create gorm db connection
	db.AutoMigrate(&Car{}, &Driver{})                            //migrate the models

	crudex.Setup() // this configuration is used by crudex to setup the controllers and scaffold behaviours
	// this is the default configuration, you can customize it by calling the methods on the returned object
	// for example: crudex.Setup().SetScaffoldRootDir("gen").SetScaffoldCreateStrategy(crudex.SCAFFOLD_ALWAYS)
	// Check the documentation for more information on the available methods, or browse the source code

	app.HTMLRender = crudex.NewRenderer() //set the renderer to the one provided by crudex

	var ctrls = []crudex.ICrudCtrl{ // create the controllers for the models
		crudex.New[Car](db).OnRouter(app.Group("cars")).ScaffoldDefaults(),
		crudex.New[Driver](db).OnRouter(app.Group("drivers")).ScaffoldDefaults(),
	}

	crudex.ScaffoldIndex(app, "gen/index.html", ctrls...) // create an index page that lists all the models
	app.Run(":8080")                                      //run the app
}

type Car struct {
	crudex.BaseModel
	//you can also customize the input type and placeholder	through the crud-input and crud-placeholder tags
	Name        string `crud-input:"text" crud-placeholder:"Enter name"`
	License     string `crud-input:"text" crud-placeholder:"Enter the license plate"`
	Description string `crud-input:"wysiwyg" crud-placeholder:"Describe it"`
	Year        int    `crud-input:"number" crud-placeholder:"Model year of the car"`
}

type Driver struct {
	crudex.BaseModel
	Name  string
	CarID uint
	Car   Car `gorm:"foreignKey:CarID"`
}
