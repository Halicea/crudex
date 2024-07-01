package crudex

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex/shared"
)

// ModelDescriptor is a struct that holds the information needed to scaffold a model.
// It is used to scaffold the templates for the given model
type ModelDescriptor struct {
	// Type is the reflect.Type of the model
	Type reflect.Type

	// Name is the name of the model
	Name string

	// TemplateFileName is the name of the file where the template will be written
	TemplateFileName string

	// Fields is a slice of reflect.StructField that represent the fields of the model that will be scaffolded
	Fields []reflect.StructField
}

// LayoutData is a struct that holds the data needed to scaffold the layout template
type LayoutData struct {
	TemplateFileName string
	Menu             []MenuItem
}

// MenuItem is a struct that holds the data needed to render a link to a model page in the layout template
type MenuItem struct {
	Title string
	Path  string
}

// ModelDescriptorConfiguration is a struct that is used to create a ModelDescriptor
// it defines the rules for the creation of the ModelDescriptor
type ModelDescriptorConfiguration struct {
	// RootDir is the root directory where the templates will be written.
	// It is used to create the TemplateFileName of the ModelDescriptor
	RootDir string
	// ModelNameSuffix is the suffix that will be added to the model name
	ModelNameSuffix string
	// TemplateNameSuffix is the suffix that will be added to the template name
	TemplateNameSuffix string
	// TemplateNamePrefix is the prefix that will be added to the template name
	TemplateNamePrefix string
}

func NewDescriptor(data interface{}, opts *ModelDescriptorConfiguration) *ModelDescriptor {
	if opts == nil {
		opts = &ModelDescriptorConfiguration{
			RootDir:            "templates",
			ModelNameSuffix:    "",
			TemplateNameSuffix: "",
		}
	}
	rootDir := opts.RootDir
	modelType := extractType(data)
	baseName := modelType.Name()
	baseNameLower := strings.ToLower(baseName)

	modelName := baseName
	templateName := baseNameLower

	if opts.ModelNameSuffix != "" {
		modelName = fmt.Sprintf("%s%s", baseName, opts.ModelNameSuffix)
	}
	if opts.TemplateNameSuffix != "" {
		templateName = fmt.Sprintf("%s-%s", baseNameLower, strings.ToLower(opts.TemplateNameSuffix))
	}
	fields := []reflect.StructField{}
	allowedTypes := []reflect.Kind{reflect.String, reflect.Int, reflect.Float64, reflect.Bool}
	for i := 0; i < modelType.NumField(); i++ {
		if !contains(allowedTypes, modelType.Field(i).Type.Kind()) {
			continue
		}
		fields = append(fields, modelType.Field(i))
	}
	return &ModelDescriptor{
		Type:             modelType,
		Name:             modelName,
		TemplateFileName: fmt.Sprintf("%s/%s.html", rootDir, templateName),
		Fields:           fields,
	}
}

func (md *ModelDescriptor) Flush(definition string, strategy ScaffoldStrategy) error {
	if !shouldScaffold(strategy, md.TemplateFileName) {
		if gin.IsDebugging() {
			fmt.Printf("Skipping scaffold of %s\n", md.TemplateFileName)
		}
		return nil
	}
	tmpl := template.Must(template.New(md.Name).
		Delims("[[", "]]").
		Funcs(config.ScaffoldMap().FuncMap()).
		Parse(definition))
	tmplFile, err := os.Create(md.TemplateFileName)
	defer tmplFile.Close()
	if err != nil {
		return err
	}
	err = tmpl.Execute(tmplFile, md)
	if err != nil {
		return err
	}
	return nil
}

func FlushAll(dst string, models ...interface{}) {
	for _, m := range models {
		GenDetailTmpl(m, dst)
		GenListTmpl(m, dst)
		GenFormTmpl(m, dst)
	}
}

func GenDetailTmpl(data interface{}, rootDir string) {
	err := NewDescriptor(data, &ModelDescriptorConfiguration{
		RootDir: rootDir,
	}).Flush(config.ScaffoldMap().Get(shared.ScaffoldTemplateDetail)(), config.ScaffoldStrategy())

	if err != nil {
		panic(err)
	}
}


func GenListTmpl(data interface{}, rootDir string) {
	err := NewDescriptor(data, &ModelDescriptorConfiguration{
		RootDir:            rootDir,
		TemplateNameSuffix: "list",
		ModelNameSuffix:    "List",
	}).Flush(config.ScaffoldMap().Get(shared.ScaffoldTemplateList)(), config.ScaffoldStrategy())

	if err != nil {
		panic(err)
	}
}

func GenFormTmpl(data interface{}, rootDir string) {
	err := NewDescriptor(data, &ModelDescriptorConfiguration{
		RootDir:            rootDir,
		TemplateNameSuffix: "form",
	}).Flush(config.ScaffoldMap().Get(shared.ScaffoldTemplateForm)(), config.ScaffoldStrategy())

	if err != nil {
		panic(err)
	}
}

func GenLayout(fileName string, controllers []ICrudCtrl) {
type TestStruct struct {
	Str  string
	Num  int32
	Html string `crud-input:"html" crud-placeholder:"Enter some HTML"`
	Date string `crud-input:"datetime"`
}
	if !shouldScaffold(config.ScaffoldStrategy(), fileName) {
		if gin.IsDebugging() {
			fmt.Printf("Skipping scaffold of %s\n", fileName)
		}
		return
	}

	data := LayoutData{
		Menu:             []MenuItem{},
		TemplateFileName: fileName,
	}

	for _, ctrl := range controllers {
		data.Menu = append(data.Menu, MenuItem{
			Title: ctrl.GetModelName(),
			Path:  ctrl.BasePath(),
		})
	}

	tmpl := template.Must(template.New(filepath.Base(fileName)).
		Delims("[[", "]]").
		Funcs(config.ScaffoldMap().FuncMap()).
		Parse(config.ScaffoldMap().Get(shared.ScaffoldTemplateLayout)()))

	tmplFile, err := os.Create(fileName)
	defer tmplFile.Close()
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(tmplFile, data)
	if err != nil {
		panic(err)
	}
}

func shouldScaffold(strategy ScaffoldStrategy, fileName string) bool {
	switch strategy {
	case SCAFFOLD_ALWAYS:
		return true
	case SCAFFOLD_NEVER:
		return false
	case SCAFFOLD_IF_NOT_EXISTS:
		_, err := os.Stat(fileName)
		return err != nil
	default:
		return false
	}
}

func extractType(data interface{}) reflect.Type {
	modelType := reflect.TypeOf(data)
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		modelType = val.Elem().Type()
	}
	return modelType
}

func contains(allows []reflect.Kind, checked reflect.Kind) bool {
	for _, a := range allows {
		if a == checked {
			return true
		}
	}
	return false
}
