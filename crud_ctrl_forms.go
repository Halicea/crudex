package crudex

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
