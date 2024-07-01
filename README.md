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
	"os"

	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := gin.New()                                             //create gin app
	db, _ := gorm.Open(sqlite.Open("sample.db"), &gorm.Config{}) //create gorm db connection
	db.AutoMigrate(&Car{}, &Driver{})                            //migrate the models

	// this configuration is used by crudex to setup the controllers and scaffold behaviours
	crudex.Setup().
		WithScaffoldStrategy(crudex.SCAFFOLD_ALWAYS).
		WithCommandLineArgs(os.Args[1:]) 
        //there are a bunch of other settings you can use, please check the documentation

	app.HTMLRender = crudex.NewRenderer()

	var controllers = []crudex.ICrudCtrl{
		crudex.New[Car](db).OnRouter(app).ScaffoldDefaults(),
		crudex.New[Driver](db).OnRouter(app).ScaffoldDefaults(),
	} // create the controllers for the models

	crudex.ScaffoldIndex(app, "gen/index.html", controllers...) // create an index page that lists all the models
	app.Run(":8080")                                            //run the app
}

type Car struct {
	crudex.BaseModel
	//you can also customize the input type and placeholder	through the crud-input and crud-placeholder tags
	Name        string `crud-input:"text" crud-placeholder:"Enter the name of the car"`
	License     string `crud-input:"text" crud-placeholder:"Enter the license plate"`
	Description string `crud-input:"wysiwyg" crud-placeholder:"Describe me"`
}

type Driver struct {
	crudex.BaseModel
	Name  string
	CarID uint
	Car   Car `gorm:"foreignKey:CarID"`
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
