package main

import (
	"github.com/42dotmk/crudex/handlers"
	"github.com/42dotmk/crudex/lib/crud"
	"github.com/42dotmk/crudex/lib/renderers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	runServer()
}

func runServer() {
	godotenv.Load() // load the .env file
	db := database.Connect()
	models.Migrate(db) // auto migrate the models and set the database

	app := setupApp() // create the gin app and setup defaults

	scaffoldAdmin(app.Group("/admin"), db)
	registerRoutes(app, db) // register the routes

	setupRenderer(app)
	app.Run()
}

func scaffoldAdmin(r crud.IRouter, db *gorm.DB) {
	controllers := []crud.ICrudCtrl{
		crud.Scaffold[models.Donation](r, db),
		crud.Scaffold[models.Repairman](r, db),
		crud.Scaffold[models.Donor](r, db),
		crud.Scaffold[models.Recipient](r, db),
		crud.Scaffold[models.EquipmentItem](r, db),
		crud.Scaffold[models.Tag](r, db),
		crud.Scaffold[models.ServicedMachine](r, db),
		crud.Scaffold[models.ServiceComment](r, db),
	}

	crud.GenLayout(controllers, "gen")
	r.GET("/", func(c *gin.Context) { renderers.Render(c, gin.H{"Path": r.BasePath()}, "layout.html") })
}

func setupApp() (app *gin.Engine) {
	app = gin.Default()
	app.Static("/assets", "./assets")
	return app
}

func setupRenderer(app *gin.Engine) {
	config := renderers.NewConfig().
		WithTemplateDirs("gen").
		WithLayout("gen/layout.html").
		WithEnableLayoutOnNonHxRequest(true).
		WithLayoutDataFunc(func(data gin.H) {
			data["Menu"] = []models.MenuItem{
				{Title: "Home", Uri: "/", IsEnabled: true, IsExternal: true},
				{Title: "Донори", Uri: "/donors", IsEnabled: true, IsExternal: false},
				{Title: "Донации", Uri: "/donations", IsEnabled: true, IsExternal: false},
				{Title: "Тим", Uri: "/team", IsEnabled: true, IsExternal: false},
				{Title: "Опрема", Uri: "/equipment", IsEnabled: true, IsExternal: false},
				{Title: "Конфигурации", Uri: "/bundles", IsEnabled: true, IsExternal: false},
			}
		})
	renderers.SetupDefaultRenderer(app, config)
}

func registerRoutes(app *gin.Engine, db *gorm.DB) {
	//V0 just create a simple method handler for a simple route (the method can be pased as an argument it does not have to be a closure)
	app.GET("/", func(c *gin.Context) { renderers.Render(c, gin.H{}, "templates/home.html") })

	// v1. create a controller for the donor model hosted at /donors implementing all its routes there
	// the positive aspect of this is that you can easily see all the routes for the donor model in one place
	// all the dependencies for the donor model are also in one place and passed to the controller explicitly
	handlers.NewDonorCtrl(db).OnRouter(app.Group("/donors")) // this method call with attach all the routes to the /donors router group

	// V2.0
	// If you need the basic CRUD operations for a model, you can just use the CrudCtrl (or inherit from it and add/override routes)
	crud.New[models.Tag](db).OnRouter(app.Group("/tags"))

	// V2.1
	// the repairman binder is a function that takes a gin context and returns a repairman or an error
	// the repairman binder is used to bind the form data to the repairman model
	crud.New[models.Repairman](db).OnRouter(app.Group("/team"))

	// v2.2 you can add the routes manually if you like
	recepients := crud.New[models.Recipient](db)
	app.GET("/recipients", recepients.List)
}
