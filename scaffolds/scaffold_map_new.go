package scaffolds

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"reflect"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex/shared"
)

type ScaffoldTemplateMap struct {
	Templates        map[string]func() *template.Template
	TemplatesFuncMap template.FuncMap
}
//go:embed scaffold_templates/**/*.html
var scaffoldTemplatesFS embed.FS

//go:embed scaffold_templates/layout/layout.html
var Layout string

//go:embed scaffold_templates/detail.html
var Detail string

//go:embed scaffold_templates/list.html
var List string

//go:embed scaffold_templates/form.html
var Form string

// New creates a new ScaffoldMap with the default templates
//
// if there is a scaffolds directory in the current working directory,
//
// it will try to read each template from the corresponding file in the directory, otherwise it will use the default template
func New() *ScaffoldMap {
	res := &ScaffoldMap{
		Templates:        make(map[string]func() string),
		TemplatesFuncMap: make(template.FuncMap),
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

func NewScaffoldTemplateMap() *ScaffoldTemplateMap {
	res := &ScaffoldTemplateMap{
		Templates:        make(map[string]func() *template.Template),
		TemplatesFuncMap: make(template.FuncMap),
	}

	res.TemplatesFuncMap["RenderInputType"] = RenderInputType
	res.Templates[shared.ScaffoldTemplateLayout.String()] = _buildScaffoldTemplate("layout", Layout, res.TemplatesFuncMap)
	return res
}

func _buildScaffoldTemplate(name string, content string, funcMap template.FuncMap) func() *template.Template {
	return func() *template.Template {
		return template.Must(template.New(name).Funcs(funcMap).Delims("[[", "]]").Parse(content))
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
// NewScaffoldTemplateRenderer creates a new multitemplate renderer using the scaffold templates
func NewScaffoldTemplateRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	templateDirs := []string{"scaffolds"} // the default template directory

	if gin.IsDebugging() {
		fmt.Fprint(gin.DefaultWriter, "Loading templates from: ", templateDirs, "\n")
	}

    // Load scaffold templates from the embedded filesystem
    err:=fs.WalkDir(scaffoldTemplatesFS, "scaffold_templates", func(path string, d fs.DirEntry, err error) error {
        if d.IsDir() {
            return nil
        }
        name := filepath.Base(path)
        if gin.IsDebugging() {
            fmt.Fprint(gin.DefaultWriter, "Loading Scaffold template: ", path, " with name ", name, "\n")
        }
        r.AddFromFS(name, scaffoldTemplatesFS, path)
        return nil
    })

    if err != nil {
        panic(err)
    }
    // Load scaffold templates from the filesystem and override the embedded ones if these exist
	for _, td := range templateDirs {
		files, err := filepath.Glob(filepath.Join(td, "*.html"))
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			name := filepath.Base(file)
			if gin.IsDebugging() {
				fmt.Fprint(gin.DefaultWriter, "Loading scaffold template: ", file, " with name ", name, "\n")
			}
			r.AddFromFiles(name, file)
		}
	}
	return r
}
