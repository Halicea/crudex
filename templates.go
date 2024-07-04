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

var SupportedScaffoldTypes = []reflect.Kind{
	reflect.String,
	reflect.Int,
	reflect.Bool,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Float32,
	reflect.Float64,
}

// ScaffoldDataModel is a struct that holds the information needed to scaffold a template for a model.
// It is used to scaffold the templates for the given model
type ScaffoldDataModel struct {
	// Type is the reflect.Type of the model
	Type reflect.Type

	// Name is the name of the model
	Name string

	// TemplateFileName is the name of the file where the template will be written
	TemplateFileName string

	// Fields is a slice of reflect.StructField that represent the fields of the model that will be scaffolded
	//
	// The fields are filtered to only include the supported types

	Fields []reflect.StructField

	// AllFields is a slice of reflect.StructField that represent all the fields of the model
	AllFields []reflect.StructField
}

// ScaffoldLayoutDataModel is a struct that holds the data needed to scaffold the layout template
type ScaffoldLayoutDataModel struct {
	TemplateFileName string
	Menu             []ScaffoldMenuItem
}

// ScaffoldMenuItem is a struct that holds the data needed to render a link to a model page in the layout template
type ScaffoldMenuItem struct {
	Title string
	Path  string
}

// ScaffoldDataModelConfigurator is a struct that is used to create a ModelDescriptor
//
// it defines the rules for the creation of the ModelDescriptor
type ScaffoldDataModelConfigurator struct {
	// RootDir is the root directory where the templates will be written.
	//
	// It is used to create the TemplateFileName of the ModelDescriptor
	RootDir string

	// ModelNameSuffix is the suffix that will be added to the model name
	ModelNameSuffix string

	// TemplateNameSuffix is the suffix that will be added to the template name
	TemplateNameSuffix string

	// TemplateNamePrefix is the prefix that will be added to the template name
	TemplateNamePrefix string

	// TemplateExtension is the extension that will be added to the template name
	TemplateExtension string
}

func NewScaffoldDataModel(data interface{}, opts *ScaffoldDataModelConfigurator) *ScaffoldDataModel {
	if opts == nil {
		panic("opts cannot be nil")
	}
	modelType := extractType(data)
	baseName := modelType.Name()
	baseNameLower := strings.ToLower(baseName)

	modelName := baseName
	templateName := baseNameLower

	if opts.ModelNameSuffix != "" {
		modelName = fmt.Sprintf("%s%s", baseName, opts.ModelNameSuffix)
	}
	if opts.TemplateNameSuffix != "" {
		templateName = fmt.Sprintf("%s%s", baseNameLower, opts.TemplateNameSuffix)
	}

	if opts.TemplateExtension != "" {
		templateName = fmt.Sprintf("%s%s", templateName, opts.TemplateExtension)
	}
	if opts.TemplateNamePrefix != "" {
		templateName = fmt.Sprintf("%s%s", opts.TemplateNamePrefix, templateName)
	}
	fields := []reflect.StructField{}
	allFields := []reflect.StructField{}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		allFields = append(allFields, field)
		if !contains(SupportedScaffoldTypes, field.Type.Kind()) {
			continue
		}
		fields = append(fields, field)
	}
	fileName := templateName
	if opts.RootDir != "" {
		fileName = fmt.Sprintf("%s/%s", opts.RootDir, templateName)
	}
	return &ScaffoldDataModel{
		TemplateFileName: fileName,
		Type:             modelType,
		Name:             modelName,
		Fields:           fields,
		AllFields:        allFields,
	}
}

func (md *ScaffoldDataModel) Flush(definition string, strategy ScaffoldStrategy) error {
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
	err := NewScaffoldDataModel(data, &ScaffoldDataModelConfigurator{
		RootDir:           rootDir,
		TemplateExtension: ".html",
	}).Flush(_scaffoldFor(shared.ScaffoldTemplateDetail), config.ScaffoldStrategy())

	if err != nil {
		panic(err)
	}
}

func GenListTmpl(data interface{}, rootDir string) {
	err := NewScaffoldDataModel(data, &ScaffoldDataModelConfigurator{
		RootDir:            rootDir,
		TemplateNameSuffix: "-list",
		ModelNameSuffix:    "List",
		TemplateExtension:  ".html",
	}).Flush(_scaffoldFor(shared.ScaffoldTemplateList), config.ScaffoldStrategy())

	if err != nil {
		panic(err)
	}
}

func GenFormTmpl(data interface{}, rootDir string) {
	err := NewScaffoldDataModel(data, &ScaffoldDataModelConfigurator{
		RootDir:            rootDir,
		TemplateNameSuffix: "-form",
		TemplateExtension:  ".html",
	}).Flush(_scaffoldFor(shared.ScaffoldTemplateForm), config.ScaffoldStrategy())

	if err != nil {
		panic(err)
	}
}

func GenLayout(fileName string, controllers []ICrudCtrl) {
	if !shouldScaffold(config.ScaffoldStrategy(), fileName) {
		if gin.IsDebugging() {
			fmt.Printf("Skipping scaffold of %s\n", fileName)
		}
		return
	}

	data := ScaffoldLayoutDataModel{
		Menu:             []ScaffoldMenuItem{},
		TemplateFileName: fileName,
	}

	for _, ctrl := range controllers {
		data.Menu = append(data.Menu, ScaffoldMenuItem{
			Title: ctrl.GetModelName(),
			Path:  ctrl.BasePath(),
		})
	}
	tmpl := template.Must(template.New(filepath.Base(fileName)).
		Delims("[[", "]]").
		Funcs(config.ScaffoldMap().FuncMap()).
		Parse(_scaffoldFor(shared.ScaffoldTemplateLayout)))

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
	case ScaffoldStrategyAlways:
		return true
	case ScaffoldStrategyNever:
		return false
	case ScaffoldStrategyIfNotExists:
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

func _scaffoldFor(kind shared.ScaffoldTemplateKind) string {
	key := kind.String()
	fn := config.ScaffoldMap().Get(key)
	return fn()
}
