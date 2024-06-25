# CRUDEX

The simplest way to build your crud APIs and views.

Crudex does not interfere with your code, it just helps you to build your App faster by providing a set of tools to build your CRUD APIs and admin interfaces.
It is based on Gin and Gorm, so you can use all the features of these libraries without any restrictions.

Crudex allows you to create and extend CRUD controllers based on the model you want to expose. It also provides a set of tools to create admin interfaces for your models.
It is flexible enough so you can customize the controllers and the admin interfaces as you need.

Crudex comes with predefined scaffold templates for the admin interfaces, but you can create your own templates and use them in your project.

At this time there is no specific configuration for the access permissions but you can handle that with some middlewarethere is no specific configuration for the access permissions but you can handle that with some middleware.

## Installation
```bash
go get 42.mk/crudex
```

## How to use
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Fruit struct {
	crudex.BaseModel //the same as the gorm.BaseModel but it adds .GetID() and .SetID(value) methods
	Name   string
	Color string
}

func main() {
	// your regular gin app
	app := gin.New()

	// create new db connection and migrate your models
	db, _ := gorm.Open(sqlite.Open("sample.db"), &gorm.Config{})
	db.AutoMigrate(&Fruit{})

	// create a configuration that your crud controllers will use
	app.HTMLRender = crudex.NewRenderer(crudex.NewConfig().
		WithTemplateDirs("gen"))
	//the Config interface has many setup options, please check the documentation for more

	// you can also scaffold an index page for your models
	crudex.ScaffoldIndex(app, "gen/index.html",
		crudex.New[Fruit](db).Scaffold(app))
	//you can also create your own controllers and use the admin to scaffold the pages
	app.Run(":8080")
}
``` 
 
## What you get
 
For every model there are six(6) routes created by default:
- `GET /model/new` Shows a form to create a new record
- `GET /model/edit/:id` Shows a form to edit a record

- `GET /model` Lists all the records (with html or json)
- `GET /model/:id` Shows a single record (with html or json)

- `POST /model` Creates a new record and redirects to the list.
  It accepts either form data or json data
- `PUT /model/:id` Updates a record and redirects to the list.
  It accepts either form data or json data
- `DELETE /model/:id` Deletes a record and redirects to the list
