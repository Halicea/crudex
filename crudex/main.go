package main
import (
	"github.com/halicea/crudex/scaffolds"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: crudex export")
		return
	}
	if os.Args[1] == "export" {
		scaffolds.ExportDefaultScaffoldTemplates()
	}
}
