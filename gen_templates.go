package crudex

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/halicea/crudex/shared"
)

type GenTemplateOptions struct {
	Name             string
	TemplateFileName string
	ScaffoldStrategy ScaffoldStrategy
}

func (self *GenTemplateOptions) ShouldScaffold() bool {
	switch self.ScaffoldStrategy {
	case ScaffoldStrategyAlways:
		return true
	case ScaffoldStrategyNever:
		return false
	case ScaffoldStrategyIfNotExists:
		_, err := os.Stat(self.TemplateFileName)
		return err != nil
	default:
		return false
	}
}

func GenAllForModels(dst string, models ...interface{}) {
	for _, m := range models {
		GenDetailTmpl(m, dst)
		GenListTmpl(m, dst)
		GenFormTmpl(m, dst)
	}
}

func GenDetailTmpl(data interface{}, rootDir string) {
	model := NewScaffoldDataModel(data, &ScaffoldDataModelConfigurator{
		RootDir:           rootDir,
		TemplateExtension: ".html",
	})
	tmplString := _scaffoldFor(shared.ScaffoldTemplateDetail)
	Must(GenTemplate(tmplString, model, &GenTemplateOptions{
		Name:             model.Name,
		TemplateFileName: model.TemplateFileName,
		ScaffoldStrategy: config.ScaffoldStrategy(),
	}))

}
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func GenListTmpl(data interface{}, rootDir string) {
	model := NewScaffoldDataModel(data, &ScaffoldDataModelConfigurator{
		RootDir:            rootDir,
		TemplateNameSuffix: "-list",
		ModelNameSuffix:    "List",
		TemplateExtension:  ".html",
	})
	scaffoldTmpl := _scaffoldFor(shared.ScaffoldTemplateList)
	Must(GenTemplate(scaffoldTmpl, model, &GenTemplateOptions{
		Name:             model.Name,
		TemplateFileName: model.TemplateFileName,
		ScaffoldStrategy: config.ScaffoldStrategy(),
	}))
}

func GenFormTmpl(data interface{}, rootDir string) {
	model := NewScaffoldDataModel(data, &ScaffoldDataModelConfigurator{
		RootDir:            rootDir,
		TemplateNameSuffix: "-form",
		TemplateExtension:  ".html",
	})
	err := GenTemplate(_scaffoldFor(shared.ScaffoldTemplateForm), model, &GenTemplateOptions{
		Name:             model.Name,
		TemplateFileName: model.TemplateFileName,
		ScaffoldStrategy: config.ScaffoldStrategy(),
	})

	if err != nil {
		panic(err)
	}
}

func GenLayout(fileName string, controllers []ICrudCtrl) {
	model := ScaffoldLayoutDataModel{
		Menu:             []ScaffoldMenuItem{},
		TemplateFileName: fileName,
	}

	for _, ctrl := range controllers {
		model.Menu = append(model.Menu, ScaffoldMenuItem{
			Title: ctrl.GetModelName(),
			Path:  ctrl.BasePath(),
		})
	}

	GenTemplate(fileName, model, &GenTemplateOptions{
		Name:             filepath.Base(fileName),
		TemplateFileName: fileName,
		ScaffoldStrategy: GetConfig().ScaffoldStrategy(),
	})
}

func GenTemplate(definition string, model interface{}, opts *GenTemplateOptions) error {
	if !opts.ShouldScaffold() {
		if gin.IsDebugging() {
			fmt.Printf("Skipping scaffold of %s\n", opts.TemplateFileName)
		}
		return nil
	}

	tmpl := _buildTemplate(opts.Name, definition)

	tmplFile, err := os.Create(opts.TemplateFileName)
	defer tmplFile.Close()
	if err != nil {
		return err
	}
	err = tmpl.Execute(tmplFile, model)
	if err != nil {
		return err
	}
	return nil
}

func _buildTemplate(name string, definition string) *template.Template {
	return template.Must(template.New(name).
		Delims("[[", "]]").
		Funcs(config.ScaffoldMap().FuncMap()).
		Parse(definition))
}

func _scaffoldFor(kind shared.ScaffoldTemplateKind) string {
	key := kind.String()
	fn := config.ScaffoldMap().Get(key)
	return fn()
}
