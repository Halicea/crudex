package crudex

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"
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

var templateFuncs = template.FuncMap{
	"RenderTypeInput": RenderTypeInput,
}
func extractType(data interface{}) reflect.Type {
	modelType := reflect.TypeOf(data)
    val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		modelType = val.Elem().Type()
	}
    return modelType
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

func (md *ModelDescriptor) Flush(definition string) error {
	tmpl := template.Must(template.New(md.Name).
		Delims("[[", "]]").
		Funcs(templateFuncs).
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
		GenDetails(m, dst)
		GenList(m, dst)
		GenForm(m, dst)
	}
}

func GenDetails(data interface{}, rootDir string) {
	err := NewDescriptor(data, &NewDescriptorConf{
		RootDir: rootDir,
	}).Flush(tmplDetail)

	if err != nil {
		panic(err)
	}
}

func GenList(data interface{}, rootDir string) {
	err := NewDescriptor(data, &NewDescriptorConf{
		RootDir:            rootDir,
		TemplateNameSuffix: "list",
		ModelNameSuffix:    "List",
	}).Flush(tmplList)

	if err != nil {
		panic(err)
	}
}

func GenForm(data interface{}, rootDir string) {
	err := NewDescriptor(data, &NewDescriptorConf{
		RootDir:            rootDir,
		TemplateNameSuffix: "form",
	}).Flush(tmplForm)

	if err != nil {
		panic(err)
	}
}

type MenuItem struct {
    Title string
    Path   string
}

func GenLayout(controllers []ICrudCtrl, rootDir string) {
    fileName := fmt.Sprintf("%s/%s", rootDir, "layout.html")
    data := struct {
        TemplateFileName string
        Menu []MenuItem
    }{
        Menu: []MenuItem{},
        TemplateFileName: fileName,
    }

    for _, ctrl := range controllers {
        data.Menu = append(data.Menu, MenuItem{
            Title: ctrl.GetModelName(),
            Path: ctrl.BasePath(),
        })
    }

	tmpl := template.Must(template.New(fileName).
		Delims("[[", "]]").
		Funcs(templateFuncs).
		Parse(tmplLayout))

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

func RenderTypeInput(field reflect.StructField) string {
	switch field.Type.Kind() {
	case reflect.String:
		return fmt.Sprintf(`<input type="text" name="%s">`, field.Name)
	case reflect.Int, reflect.Float64:
		return fmt.Sprintf(`<input type="number" name="%s">`, field.Name)
	case reflect.Bool:
		return fmt.Sprintf(`<input type="checkbox" name="%s">`, field.Name)
	}
	panic(fmt.Sprintf("unsupported type: %s for field %s", field.Type.Kind().String(), field.Name))
}

func contains(allows []reflect.Kind, checked reflect.Kind) bool {
	for _, a := range allows {
		if a == checked {
			return true
		}
	}
	return false
}
