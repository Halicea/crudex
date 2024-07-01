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
        // scaffold will generate the templates and attach this controller routes to the router
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

### Cusomization
There are three parts that you can further customize to your needs:

1. **Generated templates**
    
    The templates are exported the first time you create a new model (by default in the 'gen' directory).
    You can modify them in the way it is suitable for your use case


2. **Scaffold templates**
    
    Scaffold templates are the templates that generate the model templates used by the `CrudCtrl[T]`.

    You can export them once by invoking `crudex --export-scaffolds` in the root of your module. 
    This will create a `scaffolds` directory that will be further used to generate the CRUD templates.

    This scaffold templates can be customized to generate new CRUD templates with the look, feel and functionality suitable to You.

3. **Controller**

    You can create your own controller that builds on top CrudCtrl[T]

    This way you can override the `List`, `Details`, `Form`, `List` handlers and add your own.

    Behind the scenes crudex uses `gin` handlers, so you can build any route without additional need to learn something new.
    
## Wishlist

- [X] **[P1]** Start with tests
- [X] **[P1]** Add more customization options
- [X] **[P1]** Initial README

- [ ] **[P3]** Use source generators to scaffold the templates (through `go generate`)
- [ ] **[P3]** Create separate package for the template scaffolding and leave just the controllers in this package
- [ ] **[P3]** Add more documentation 
- [ ] **[P3]** Fully document the public methods, interfaces and structs
- [ ] **[P2]** Add more tests
- [ ] **[P2]** Allow the possibility for different UI packages to be glued to it
     - For example:
        - generate UI templates with daisyUI and React
        - generate UI templates with HTMX and tailwind
        - e.t.c
