package crudex
// {{/* generated file: [[.TemplateFileName]] */}}
// <section>
//     [[$modelName := .Name]]
//     <h1>[[$modelName]]</h1>
//     <form
//         {{if .ID}}hx-post="{{.Path}}"{{else}}hx-put="{{.Path}}/new"{{end}}
//         hx-target="#main">[[range .Fields]]
//         <div>
//             <label for="[[.Name]]">[[.Name]]</label>
//             [[RenderInputType $modelName .]]
//         </div>[[end]]
//         <button type="submit">Submit</button>
//     </form>
// <section>
