package crudex

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FormBinder is a function that binds the form data to a model
type FormBinder[T IModel] func(c *gin.Context, out *T) error

// CrudCtrl is a controller that implements the basic CRUD operations for a model with a gorm backend
type CrudCtrl[T IModel] struct {
	Db     *gorm.DB
	Router IRouter
	Config IConfig

	ModelName    string
	ModelKeyName string
	FormBinder   FormBinder[T]
}

// Returns the Name of the model
func (self *CrudCtrl[T]) GetModelName() string {
	return self.ModelName
}

// BasePath returns the base path of the controller
func (self *CrudCtrl[T]) BasePath() string {
	return self.Router.BasePath()
}

// New creates a new CRUD controller for the provided model
//
// It uses the default configuration
// See `NewWithOptions` for more control over the configuration
func New[T IModel]() *CrudCtrl[T] {
	conf := GetConfig()
	db := conf.DefaultDb()
	modelType := extractType(*new(T))
	ctrlGroup := conf.DefaultRouter().Group(fmt.Sprintf("/%s", strings.ToLower(modelType.Name())))
	return NewWithOptions[T](db, ctrlGroup, conf)
}

// New creates a new CRUD controller for the provided model
//
// It uses the provided configuration to define its behaviour
func NewWithOptions[T IModel](db *gorm.DB, router IRouter, conf IConfig) *CrudCtrl[T] {
	var name = fmt.Sprintf("%T", *new(T))
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[1]
	}
	res := &CrudCtrl[T]{
		FormBinder: DefaultFormHandler[T], // default form handler is used if none is provided
		ModelName:  name,
		Db:         db,
		Config:     conf,
		Router:     router,
	}
	if router != nil {
		res.OnRouter(router)
	} else {
		slog.Warn("Router is nil for %s", slog.Any("name", name))
	}
	if db == nil {
		slog.Warn("DB is nil for model", slog.Any("name", name))
	}
	if conf.AutoScaffold() {
		res.ScaffoldDefaults()
	}
	return res
}

// OnRouter attaches the CRUD routes to the provided router
func (self *CrudCtrl[T]) OnRouter(r IRouter) *CrudCtrl[T] {
	self.Router = r
	if r != nil {
		self.EnableRoutes(r)
	}
	return self
}

func (self *CrudCtrl[T]) EnableRoutes(r IRouter) *CrudCtrl[T] {
	if r == nil {
		panic("Router is nil, cannot enable routes")
	}
	r.GET("/", self.List)

	if self.Config.HasUI() {
		r.GET("/new", self.Form)
		r.GET("/:id/edit", self.Form)
	}

	r.PUT("/new", self.Upsert)

	r.GET("/:id", self.Details)
	r.POST("/:id", self.Upsert)
	r.DELETE("/:id", self.Delete)
	return self
}

func (self *CrudCtrl[T]) ScaffoldDefaults() *CrudCtrl[T] {
	model := *new(T)
	rootDir := self.Config.ScaffoldRootDir()

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		if os.MkdirAll(rootDir, 0755) != nil {
			panic("Failed to create directory")
		}
	}
	GenListTmpl(model, rootDir)
	GenDetailTmpl(model, rootDir)
	GenFormTmpl(model, rootDir)
	return self
}

// Creates and Flushes custom scaffold template
func (self *CrudCtrl[T]) Scaffold(scaffoldTmpl string, conf *ScaffoldDataModelConfigurator) *CrudCtrl[T] {
	model := *new(T)
	rootDir := self.Config.ScaffoldRootDir()
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		if os.MkdirAll(rootDir, 0755) != nil {
			panic("Failed to create directory")
		}
	}
	scafoldModel := NewScaffoldDataModel(model, conf)
	err := GenTemplate(scaffoldTmpl, scafoldModel, &GenTemplateOptions{
		Name:             scafoldModel.Name,
		TemplateFileName: scafoldModel.TemplateFileName,
		ScaffoldStrategy: GetConfig().ScaffoldStrategy(),
	})
	if err != nil {
		panic(err)
	}
	return self
}

// WithFormBinder sets the form binder for the controller to be used when binding form data to the model on the POST and PUT requests.
// It is used in the Upsert method of the controller
// If not set, the default form binder is used which assumes that the form field names are the same as the model field names(case sensitive)
func (self *CrudCtrl[T]) WithFormBinder(handler FormBinder[T]) *CrudCtrl[T] {
	self.FormBinder = handler
	return self
}

// List is a handler that lists all the items of the model
// it is a GET request
// !Requres the template to be named as modelName-list.html where the modelName is lowercased model name
func (self *CrudCtrl[T]) List(c *gin.Context) {
	var items []T
	filter, error := NewSearchArgsFromQuery(c)
	if error != nil {
		err := c.Error(error)
		c.String(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	self.Db.Find(&items).Limit(filter.Limit).Offset((filter.Page - 1) * filter.Limit)

	Respond(c,
		gin.H{fmt.Sprintf("%sList", self.ModelName): &items, "Path": self.Router.BasePath()},
		fmt.Sprintf("%s-list.html", strings.ToLower(self.ModelName)))
}

// Details is a handler that shows the details of a single item of the model
// it is a GET request
// !Requires the template to be named as modelName-details.html where the modelName is lowercased model name
func (self *CrudCtrl[T]) Details(c *gin.Context) {
	template := fmt.Sprintf("%s.html", strings.ToLower(self.ModelName))
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err == nil {
		var item T
		self.Db.First(&item, id)
		Respond(c, gin.H{self.ModelName: item, "Path": fmt.Sprintf("%s/%s", self.Router.BasePath(), idStr)}, template)
	} else {
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid ID for %s: %d", self.ModelName, id))
	}
}

// Form is a handler that shows the form for editing an item of the model
// it is a GET request
// !Requires the template to be named as modelName-edit.html where the modelName is lowercased model name
func (self *CrudCtrl[T]) Form(c *gin.Context) {
	template := fmt.Sprintf("%s-form.html", strings.ToLower(self.ModelName))
	idStr := c.Param("id")
	if idStr == "" {
		Respond(c, gin.H{"Path": self.Router.BasePath()}, template)
	} else {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err == nil {
			var item T
			self.Db.First(&item, id)
			Respond(c, gin.H{self.ModelName: item, "Path": fmt.Sprintf("%s/%s", self.Router.BasePath(), idStr)}, template)
		} else {
			c.String(http.StatusBadRequest, fmt.Sprintf("Invalid ID for %s: %d", self.ModelName, id))
		}
	}
}

// Upsert is a handler that saves an item of the model
// it is a POST or PUT request depending on the presence of the id parameter
// !Requires the form fields to be named as the model field names(case sensitive)
// It redirects to the details page of the saved item
func (self *CrudCtrl[T]) Upsert(c *gin.Context) {
	var item T
	if err := self.FormBinder(c, &item); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}
	idStr := c.Param("id")
	isNew := idStr == ""
	if !isNew {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Invalid ID: %d", id)
			c.String(http.StatusBadRequest, msg)
			c.Abort()
			return
		}
		item.SetID(uint(id))
	}
	res := self.Db.Save(&item)

	if res.Error != nil {
		c.String(http.StatusBadRequest, res.Error.Error())
		c.Abort()
		return
	}
	c.Header("HX-Redirect", fmt.Sprintf("%s/%d", self.BasePath(), item.GetID()))
	c.String(http.StatusOK, "Saved")
	c.Abort()
}

// Delete is a handler that deletes an item of the model
// it is a DELETE request
func (self *CrudCtrl[T]) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Invalid ID: %d", id)
		c.String(http.StatusBadRequest, msg)
		c.Abort()
		return
	}
	var item T
	self.Db.Delete(&item, id)
	c.Header("HX-Redirect", self.Router.BasePath())
	c.String(http.StatusOK, "Deleted")
	c.Abort()
}
