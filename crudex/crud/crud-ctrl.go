package crud

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/halicea/crudex/renderers"
)

type IRouter interface {
	gin.IRoutes
	Group(string, ...gin.HandlerFunc) *gin.RouterGroup
	BasePath() string
}

type IModel interface {
	GetID() uint
	SetID(id uint)
}

// FormBinder is a function that binds the form data to a model
type FormBinder[T IModel] func(c *gin.Context, out *T) error

// ICrudCtrl is an interface that defines the basic CRUD operations for a model
type ICrudCtrl interface {
	BasePath() string
	GetModelName() string

	List(c *gin.Context)
	Details(c *gin.Context)
	Form(c *gin.Context)
	Upsert(c *gin.Context)
	Delete(c *gin.Context)
}

// CrudCtrl is a controller that implements the basic CRUD operations for a model with a gorm backend
type CrudCtrl[T IModel] struct {
	Db           *gorm.DB
	Router       IRouter
	ModelName    string
	ModelKeyName string
	FormBinder   FormBinder[T]
}

func (self *CrudCtrl[T]) GetModelName() string {
	return self.ModelName
}

func (self *CrudCtrl[T]) BasePath() string {
	return self.Router.BasePath()
}

// Scaffold creates the CRUD controllers for the provided model, generates templates, and attaches them to the provided router
func Scaffold[T IModel](r IRouter, db *gorm.DB) ICrudCtrl {
	model := *new(T)
	dst := "gen"

	GenDetails(model, dst)
	GenList(model, dst)
	GenForm(model, dst)

	//create routes for the model
	modelType := extractType(model)
	ctrlGroup := r.Group(fmt.Sprintf("/%s", strings.ToLower(modelType.Name())))
	ctrl := New[T](db).OnRouter(ctrlGroup)
	return ctrl
}

// New creates a new CRUD controller for the provided model
func New[T IModel](db *gorm.DB, ) *CrudCtrl[T] {
	var name = fmt.Sprintf("%T", *new(T))
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[1]
	}
	name = strings.ToLower(name)

	return &CrudCtrl[T]{
		FormBinder: DefaultFormHandler[T], // default form handler is used if none is provided
		ModelName:  name,
		Db:         db,
	}
}

// WithModelKeyName sets the name of the primaryKey of the model, we will search by this key
// func (self *CrudCtrl[T]) WithModelKeyName(key string) *CrudCtrl[T]{
//     if key == "" {
//         panic("Invalid key name")
//     }
//     self.ModelKeyName = key
//     return self
// }

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
	self.Db.Find(&items)
	renderers.Render(c,
		gin.H{fmt.Sprintf("%sList", self.ModelName): items, "Path": c.Request.URL.Path},
		fmt.Sprintf("%s-list.html", self.ModelName))
}

// Details is a handler that shows the details of a single item of the model
// it is a GET request
// !Requires the template to be named as modelName-details.html where the modelName is lowercased model name
func (self *CrudCtrl[T]) Details(c *gin.Context) {
	template := fmt.Sprintf("%s.html", self.ModelName)
	modelName := fmt.Sprintf("%s", self.ModelName)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err == nil {
		var item T
		self.Db.First(&item, id)
		renderers.Render(c, gin.H{modelName: item, "Path": c.Request.URL.Path}, template)
	} else {
		c.Error(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Invalid ID for %s: %d", self.ModelName, id))
	}
}

// Form is a handler that shows the form for editing an item of the model
// it is a GET request
// !Requires the template to be named as modelName-edit.html where the modelName is lowercased model name
func (self *CrudCtrl[T]) Form(c *gin.Context) {
	template := fmt.Sprintf("%s-form.html", self.ModelName)
	modelName := fmt.Sprintf("%s", self.ModelName)
	idStr := c.Param("id")
	if idStr == "" {
		renderers.Render(c, gin.H{"Path": c.Request.URL.Path}, template)
		return
	} else {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err == nil {
			var item T
			self.Db.First(&item, idStr)
			renderers.Render(c, gin.H{modelName: item, "Path": c.Request.URL.Path}, template)
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
	c.Header("HX-Redirect", fmt.Sprintf("%s%s", c.Request.URL.Path, self.ModelName))
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
