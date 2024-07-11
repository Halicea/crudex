package crudex

import (
	"embed"
	"testing"

	"github.com/halicea/crudex"
	"github.com/halicea/crudex/scaffolds"
)

//go:embed *.html
var templatesFS embed.FS

type TestModel struct {
	crudex.BaseModel
	Name string
}
type FSScaffoldMap struct{ scaffolds.ScaffoldMap }

func TestGenDetailTmpl(t *testing.T) {

	// Output:
	// <section>
	//     <h1>Model</h1>
	//     <div>
	//         <div>
	//             <label for="ID">ID</label>
	//             <div>{{.Model.ID}}</div>
	//         </div>
	//         <div>
	//             <label for="Name">Name</label>
	//             <div>{{.Model.Name}}</div>
	//         </div>
	//     </div>
	// </section>
}

// <section>
//     [[$modelName := .Name]]
//     <h1>[[$modelName]]</h1>
//     <div>
//         <div>
//             <label for="ID">ID</label>
//             <div>{{.[[$modelName]].ID}}</div>
//         </div>
//         [[range .Fields]]
//         <div>
//             <label for="[[.Name]]">[[.Name]]</label>
//             <div>{{.[[$modelName]].[[.Name]]}}</div>
//         </div>
//     [[end]]</div>
// </section>
