package scaffolds

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"html/template"
	"github.com/halicea/crudex/shared"
)


type ScaffoldMap struct {
	Templates        map[string]func() string
	TemplatesFuncMap template.FuncMap
}


func (self *ScaffoldMap) FuncMap() template.FuncMap {
	return self.TemplatesFuncMap
}

func (self *ScaffoldMap) Get(name string) func() string {
	return self.Templates[name]
}

func (self *ScaffoldMap) GetString(name string) string {
	return self.Templates[name]()
}

// All returns all the scaffold templates in the map
func (self *ScaffoldMap) All() map[string]func() string {
	return self.Templates
}

// WithFuncMap sets the FuncMap that will be passed to the scaffold templates
func (self *ScaffoldMap) WithFuncMap(funcMap template.FuncMap) *ScaffoldMap {
	self.TemplatesFuncMap = funcMap
	return self
}

// Set sets the scaffold template that generates the content for the given name
//
// name is the name of the template
// `function` is a function that returns the content of the template
func (self *ScaffoldMap) Set(
	name string,
	function func() string) *ScaffoldMap {
	self.Templates[name] = function
	return self
}

// SetString sets the scaffold template that generates the content for the given name
func (self *ScaffoldMap) SetString(name string, content string) *ScaffoldMap {
	self.Templates[name] = func() string { return content }
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
func (self *ScaffoldMap) PrintAll(w io.Writer) {
	_write(w, "Scaffold Templates:\n")
	for name, content := range self.All() {
		_write(w, fmt.Sprintf("Template: %s\n", name))
		_write(w, content())
		_write(w, "============================\n")
	}
}



func (self *ScaffoldMap) String() string {
	return fmt.Sprintf(`
        Templates: %v,
        TemplatesFuncMap: %v,`,
		self.Templates,
		self.TemplatesFuncMap)
}

//_write writes the content to the writer and panics if there is an error
func _write(w io.Writer, content string) {
	_, err := w.Write([]byte(content))
	if err != nil {
		panic(err)
	}
}
