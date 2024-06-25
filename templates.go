package crudex

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

const (
	SCAFFOLD_ALWAYS = iota
	SCAFFOLD_IF_NOT_EXISTS
	SCAFFOLD_NEVER
)

type ModelDescriptor struct {
	Type             reflect.Type
	Name             string
	TemplateFileName string
	Fields           []reflect.StructField
}

type NewDescriptorConf struct {
	RootDir            string
	ModelNameSuffix    string
	TemplateNameSuffix string

	//used to prefix the template name (this allows to override some templates)
	TemplateNamePrefix string
}

func NewDescriptor(data interface{}, opts *NewDescriptorConf) *ModelDescriptor {
	if opts == nil {
		opts = &NewDescriptorConf{
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
		Funcs(config.ScaffoldFuncs()).
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
		GenDetailsTmpl(m, dst)
		GenListTmpl(m, dst)
		GenFormTmpl(m, dst)
	}
}

func GenDetailsTmpl(data interface{}, rootDir string) {
	err := NewDescriptor(data, &NewDescriptorConf{
		RootDir: rootDir,
	}).Flush(config.DetailScaffold(), config.ScaffoldStrategy())

	if err != nil {
		panic(err)
	}
}

func GenListTmpl(data interface{}, rootDir string) {
	err := NewDescriptor(data, &NewDescriptorConf{
		RootDir:            rootDir,
		TemplateNameSuffix: "list",
		ModelNameSuffix:    "List",
	}).Flush(config.ListScaffold(), config.ScaffoldStrategy())

	if err != nil {
		panic(err)
	}
}

func GenFormTmpl(data interface{}, rootDir string) {
	err := NewDescriptor(data, &NewDescriptorConf{
		RootDir:            rootDir,
		TemplateNameSuffix: "form",
	}).Flush(config.FormScaffold(), config.ScaffoldStrategy())

	if err != nil {
		panic(err)
	}
}

type MenuItem struct {
	Title string
	Path  string
}

func GenLayout(fileName string, controllers []ICrudCtrl) {
	if !shouldScaffold(config.ScaffoldStrategy(), fileName) {
		if gin.IsDebugging() {
			fmt.Printf("Skipping scaffold of %s\n", fileName)
		}
		return
	}

	data := struct {
		TemplateFileName string
		Menu             []MenuItem
	}{
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
		Funcs(config.ScaffoldFuncs()).
		Parse(config.LayoutScaffold()))

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
