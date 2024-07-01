package crudex

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FormBinder is a function that binds the form data to a model
type FormBinder[T IModel] func(c *gin.Context, out *T) error

// CrudCtrl is a controller that implements the basic CRUD operations for a model with a gorm backend
type CrudCtrl[T IModel] struct {
	Db           *gorm.DB
	Router       IRouter
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
func New[T IModel](db *gorm.DB) *CrudCtrl[T] {
	var name = fmt.Sprintf("%T", *new(T))
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[1]
	}
	return &CrudCtrl[T]{
		FormBinder: DefaultFormHandler[T], // default form handler is used if none is provided
		ModelName:  name,
		Db:         db,
	}
}


func (self *CrudCtrl[T]) ScaffoldDefaults() *CrudCtrl[T] {
	model := *new(T)
	rootDir := config.ScaffoldRootDir()
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

// OnRouter attaches the CRUD routes to the provided router
func (self *CrudCtrl[T]) OnRouter(r IRouter) *CrudCtrl[T] {
	self.Router = r

	r.GET("/", self.List)

	r.GET("/new", self.Form)
	r.PUT("/new", self.Upsert)

	r.GET("/:id", self.Details)
	r.GET("/:id/edit", self.Form)
	r.POST("/:id", self.Upsert)

	r.DELETE("/:id", self.Delete)

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
		c.Error(error)
		c.String(http.StatusBadRequest, error.Error())
		c.Abort()
		return
	}
	self.Db.Find(&items).Limit(filter.Limit).Offset((filter.Page - 1) * filter.Limit)
	Render(c,
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
		Render(c, gin.H{self.ModelName: item, "Path": fmt.Sprintf("%s/%s", self.Router.BasePath(), idStr)}, template)
	} else {
		c.Error(err)
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
		Render(c, gin.H{"Path": self.Router.BasePath()}, template)
		return
	} else {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err == nil {
			var item T
			self.Db.First(&item, id)
			Render(c, gin.H{self.ModelName: item, "Path": fmt.Sprintf("%s/%s", self.Router.BasePath(), idStr)}, template)
		} else {
			c.Error(err)
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
		c.Error(err)
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
		c.Error(res.Error)
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

func BindForm[T any](r *http.Request, out *T) error {
	// Parse the form data from the request.
	if err := r.ParseForm(); err != nil {
		return err
	}
	// Reflect on the struct to set values.
	val := reflect.ValueOf(out).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		// Get form value for the field name.
		formValue := r.FormValue(fieldType.Name)
		// Check if the field can be set and if the form value is not empty.
		if field.CanSet() && formValue != "" {
			// Convert form values to the appropriate field types.
			// This example assumes all fields are strings for simplicity.
			// You might need to convert this based on the field type.
			field.SetString(formValue)
		}
	}
	return nil
}

// ScaffoldIndex creates a simple index page that lists all the controllers
func ScaffoldIndex(r IRouter, fileName string, controllers ...ICrudCtrl) gin.IRoutes {
	GenLayout(fileName, controllers)
	return r.GET("/", func(c *gin.Context) {
		RenderWithConfig(c, gin.H{"Path": r.BasePath()}, filepath.Base(fileName), NewConfig().WithEnableLayoutOnNonHxRequest(false))
	})
}

// DefaultFormHandler is a default form binder that binds the form data to a model using the form field names as the model field names
func DefaultFormHandler[T IModel](c *gin.Context, out *T) error {
	if err := c.Request.ParseForm(); err != nil {
		return err
	}
	val := reflect.ValueOf(out).Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		formValue := c.PostForm(fieldType.Name)

		// we only set the field if we are able to do so and the form value is not empty
		if field.CanSet() && formValue != "" {
			switch fieldType.Type.Kind() {
			case reflect.Uint:
				val, err := strconv.ParseUint(formValue, 10, 64)
				if err != nil {
					return err
				}
				field.SetUint(val)
			case reflect.String:
				field.SetString(formValue)
			case reflect.Int:
				val, err := strconv.ParseInt(formValue, 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(val)
			case reflect.Float32:
				val, err := strconv.ParseFloat(formValue, 32)
				if err != nil {
					return err
				}
				field.SetFloat(val)
			case reflect.Float64:
				val, err := strconv.ParseFloat(formValue, 64)
				if err != nil {
					return err
				}
				field.SetFloat(val)
			case reflect.Bool:
				if formValue == "true" || formValue == "checked" || formValue == "1" {
					field.SetBool(true)
				} else {
					field.SetBool(false)
				}
			default:
				return fmt.Errorf("Unsupported type: %s", fieldType.Type.Kind())
			}
		}
	}
	return nil
}
