package scaffolds

import (
	_ "embed"
	"fmt"
	"os"
	"reflect"
	"text/template"

	"github.com/halicea/crudex/shared"
)

//go:embed scaffold_templates/layout.html
var Layout string

//go:embed scaffold_templates/detail.html
var Detail string

//go:embed scaffold_templates/list.html
var List string

//go:embed scaffold_templates/form.html
var Form string

type ScaffoldMap struct {
	templates map[string]func() string
	funcMap   template.FuncMap
}

func (self *ScaffoldMap) String() string {
	return fmt.Sprintf(`
        Templates: %v,
        FuncMap: %v,`,
		self.templates,
		self.funcMap)

}

func (self *ScaffoldMap) FuncMap() template.FuncMap {
	return self.funcMap
}

func (self *ScaffoldMap) Get(name string) func() string {
	return self.templates[name]
}

func (self *ScaffoldMap) GetString(name string) string {
	return self.templates[name]()
}

// All returns all the scaffold templates in the map
func (self *ScaffoldMap) All() map[string]func() string {
	return self.templates
}

// WithFuncMap sets the FuncMap that will be passed to the scaffold templates
func (self *ScaffoldMap) WithFuncMap(funcMap template.FuncMap) *ScaffoldMap {
	self.funcMap = funcMap
	return self
}

// Set sets the scaffold template that generates the content for the given name
//
// name is the name of the template
// `function` is a function that returns the content of the template
func (self *ScaffoldMap) Set(
	name string,
	function func() string) *ScaffoldMap {
	self.templates[name] = function
	return self
}

// SetString sets the scaffold template that generates the content for the given name
func (self *ScaffoldMap) SetString(name string, content string) *ScaffoldMap {
	self.templates[name] = func() string { return content }
	return self
}

// WithListScaffold sets the scaffold template that generages the list[T] template
func (self *ScaffoldMap) WithListScaffold(value func() string) *ScaffoldMap {
	return self.Set(shared.ScaffoldTemplateList.String(), value)
}

// WithDetailScaffold sets the scaffold template function that generages the detail[T] template
func (self *ScaffoldMap) WithDetailScaffold(value func() string) *ScaffoldMap {
	return self.Set(shared.ScaffoldTemplateDetail.String(), value)
}

// WithFormScaffold sets the scaffold template function that generages the form[T] template
func (self *ScaffoldMap) WithFormScaffold(value func() string) *ScaffoldMap {
	return self.Set(shared.ScaffoldTemplateForm.String(), value)
}

// WithLayoutScaffold sets the scaffold template function that generates the layout template used for listing all the models
func (self *ScaffoldMap) WithLayoutScaffold(value func() string) *ScaffoldMap {
	return self.Set(shared.ScaffoldTemplateLayout.String(), value)
}

// New creates a new ScaffoldMap with the default templates
//
// if there is a scaffolds directory in the current working directory,
//
// it will try to read each template from the corresponding file in the directory, otherwise it will use the default template
func New() *ScaffoldMap {
	res := ScaffoldMap{
		templates: make(map[string]func() string),
		funcMap:   make(template.FuncMap),
	}
	return res.
		Set(shared.ScaffoldTemplateLayout.String(), func() string { return ReadContentsOrDefault("scaffolds/layout.html", Layout, true) }).
		Set(shared.ScaffoldTemplateList.String(), func() string { return ReadContentsOrDefault("scaffolds/list.html", List, true) }).
		Set(shared.ScaffoldTemplateDetail.String(), func() string { return ReadContentsOrDefault("scaffolds/detail.html", Detail, true) }).
		Set(shared.ScaffoldTemplateForm.String(), func() string { return ReadContentsOrDefault("scaffolds/form.html", Form, true) }).
		WithFuncMap(template.FuncMap{
			"RenderInputType": RenderInputType,
		})
}

// Export writes the scaffold templates to the scaffolds directory in the current working directory.
//
// If the directory does not exist, it will be created.
//
// If the forceIfExists parameter is set to true, the files will be overwritten if they already exist.
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
	updates := 0
	for name, content := range self.All() {
		if writeIfNeeded(fmt.Sprintf("scaffolds/%s.html", name), content()) {
			updates++
		}
	}
	return nil
}

// PrintAll prints all the scaffold templates to the console
func (self *ScaffoldMap) PrintAll() {
	println("Scaffold Templates:")
	for name, content := range self.All() {
		println(fmt.Sprintf("Template: %s", name))
		println(content())
		println("============================")
	}
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
		case shared.INPUT_MARKDOWN.String(), shared.INPUT_HTML.String(), shared.INPUT_WYSIWYG.String(), shared.INPUT_TEXT.String():
			return fmt.Sprintf(`<input type="textarea" name="%s"%s>{{.%s.%s}}</input>`, field.Name, placeholder, modelName, field.Name)
		case shared.INPUT_DATETIME.String():
			return fmt.Sprintf(`<input type="datetime" name="%s"%s>{{.%s.%s}}</input>`, field.Name, placeholder, modelName, field.Name)
		default:
			panic(fmt.Sprintf("Unsupported input type '%s' specified for %s/%s", inpTag, modelName, field.Name))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Int, reflect.Float64, reflect.Int32, reflect.Int16, reflect.Int64, reflect.Int8:
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
