package scaffolds

import (
	_ "embed"
	"os"
)

//go:embed layout.html
var Layout string

//go:embed detail.html
var Detail string

//go:embed list.html
var List string

//go:embed form.html
var Form string

func ExportDefaultScaffoldTemplates() error {
	if _, err := os.Stat("scaffolds"); os.IsNotExist(err) {
		err := os.MkdirAll("scaffolds", 0755)
		if err != nil {
			return err
		}
	}

	err := os.WriteFile("scaffolds/layout.html", []byte(Layout), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile("scaffolds/detail.html", []byte(Detail), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile("scaffolds/list.html", []byte(List), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile("scaffolds/form.html", []byte(Form), 0644)
	if err != nil {
		return err
	}
	return nil
}

func PrintScaffoldTemplates() {
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
