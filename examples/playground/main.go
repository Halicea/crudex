package main

import (
	"github.com/halicea/crudex"
)

type Flower struct {
	crudex.BaseModel
	Name  string
	Color string
}

func main() {
	const tmpl = `#[[.Name]] [[range .Fields]]
    [[.Name]] [[.Type]][[end]]`

	type Conf = crudex.ScaffoldDataModelConfigurator
	crudex.New[Flower]().Scaffold(tmpl, &Conf{TemplateExtension: ".model"})
}
