package scaffolds

import (
	_ "embed"
	"fmt"
	"github.com/halicea/crudex/shared"
	"os"
	"reflect"
	"text/template"
)

//go:embed layout.html
var Layout string

//go:embed detail.html
var Detail string

//go:embed list.html
var List string

//go:embed form.html
var Form string

type ScaffoldMap struct {
	templates map[string]func() string
	funcMap   template.FuncMap
}

func (self *ScaffoldMap) FuncMap() template.FuncMap {
	return self.funcMap
}

func (self *ScaffoldMap) Get(name string) func() string {
	return self.templates[name]
}

func (self *ScaffoldMap) All() map[string]func() string {
	return self.templates
}

func (self *ScaffoldMap) WithFuncMap(funcMap template.FuncMap) *ScaffoldMap {
	self.funcMap = funcMap
	return self
}

func (self *ScaffoldMap) Set(name string, function func() string) *ScaffoldMap {
	self.templates[name] = function
	return self
}

// WithListScaffold sets the scaffold template that generages the list[T] template
func (self *ScaffoldMap) WithListScaffold(value func() string) *ScaffoldMap {
	return self.Set(shared.ScaffoldTemplateList, value)
}

// WithDetailScaffold sets the scaffold template function that generages the detail[T] template
func (self *ScaffoldMap) WithDetailScaffold(value func() string) *ScaffoldMap {
	return self.Set(shared.ScaffoldTemplateDetail, value)
}

// WithFormScaffold sets the scaffold template function that generages the form[T] template
func (self *ScaffoldMap) WithFormScaffold(value func() string) *ScaffoldMap {
	return self.Set(shared.ScaffoldTemplateForm, value)
}

// WithLayoutScaffold sets the scaffold template function that generates the layout template used for listing all the models
func (self *ScaffoldMap) WithLayoutScaffold(value func() string) *ScaffoldMap {
	return self.Set(shared.ScaffoldTemplateLayout, value)
}

func New() *ScaffoldMap {
	res := ScaffoldMap{
		templates: make(map[string]func() string),
		funcMap:   make(template.FuncMap),
	}
	return res.
		Set(shared.ScaffoldTemplateLayout, func() string { return readContentsOrDefault("scaffolds/layout.html", Layout, true) }).
		Set(shared.ScaffoldTemplateList, func() string { return readContentsOrDefault("scaffolds/list.html", List, true) }).
		Set(shared.ScaffoldTemplateDetail, func() string { return readContentsOrDefault("scaffolds/detail.html", Detail, true) }).
		Set(shared.ScaffoldTemplateForm, func() string { return readContentsOrDefault("scaffolds/form.html", Form, true) }).
		WithFuncMap(template.FuncMap{
			"RenderTypeInput": RenderInputType,
		})
}

func (self *ScaffoldMap) Export(forceIfExists bool) error {
	if _, err := os.Stat("scaffolds"); os.IsNotExist(err) {
		err := os.MkdirAll("scaffolds", 0755)
		if err != nil {
			return err
		}
	}
	writeIfNeeded := func(filename string, content string) bool {
		if !forceIfExists {
			if _, err := os.Stat(filename); err == nil {
				return false
			} else {
				if !os.IsNotExist(err) { //something else went wrong
					panic(err)
				}
			}
		}
		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			panic(err)
		}
		return false
	}

	writeIfNeeded("scaffolds/layout.html", Layout)
	writeIfNeeded("scaffolds/detail.html", Detail)
	writeIfNeeded("scaffolds/list.html", List)
	writeIfNeeded("scaffolds/form.html", Form)
	return nil
}

func PrintScaffolds() {
	println("Scaffolds")
	println("Layout")
	println(Layout)
	println("================================")
	println("Detail")
	println(Detail)
	println("================================")
	println("List")
	println(List)
	println("================================")
	println("Form")
	println(Form)
	println("================================")
}

// RenderInputType is a helper function that renders an input based on the type of the field.
//
// This function is part of the default FuncMap that is passed to the scaffold templates.
// It is used in the form template to render the input fields for the model
func RenderInputType(modelName string, field reflect.StructField) string {
	inpTag := field.Tag.Get("crud-input")
	placeholder := field.Tag.Get("crud-placeholder")
	if placeholder != "" {
		placeholder = fmt.Sprintf(" placeholder=\"%s\"", placeholder)
	}
	switch field.Type.Kind() {
	case reflect.String:
		switch inpTag {
		case "":
			return fmt.Sprintf(`<input type="text" name="%s"%s>{{.%s.%s}}</input>`, field.Name, placeholder, modelName, field.Name)
		case shared.INPUT_MARKDOWN, shared.INPUT_HTML, shared.INPUT_WYSIWYG, shared.INPUT_TEXT:
			return fmt.Sprintf(`<input type="textarea" name="%s"%s>{{.%s.%s}}</input>`, field.Name, placeholder, modelName, field.Name)
		case shared.INPUT_DATETIME:
			return fmt.Sprintf(`<input type="datetime" name="%s"%s>{{.%s.%s}}</input>`, field.Name, placeholder, modelName, field.Name)
		default:
			panic(fmt.Sprintf("Unsupported input type '%s' specified for %s/%s", inpTag, modelName, field.Name))
		}
	case reflect.Int, reflect.Float64, reflect.Int32, reflect.Int16, reflect.Int64, reflect.Int8:
		//TODO: need to make this work better
		switch inpTag {
		case "":
			return fmt.Sprintf(`<input type="number" name="%s"%s>{{.%s.%s}}</input>`, field.Name, placeholder, modelName, field.Name)
		default:
			return fmt.Sprintf(`<input type="number" name="%s"%s>{{.%s.%s}}</input>`, field.Name, placeholder, modelName, field.Name)
		}

	case reflect.Bool:
		return fmt.Sprintf(`<input type="checkbox" name="%s"%s>{{.%s.%s}}</input>`, field.Name, placeholder, modelName, field.Name)
	}

	panic(fmt.Sprintf("unsupported type: %s for field %s", field.Type.Kind().String(), field.Name))
}
