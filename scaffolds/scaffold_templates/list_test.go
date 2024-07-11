package crudex
// {{/* generated file: [[.TemplateFileName]] */}}
// <section>
//     [[$modelName := .Name]]
//     <h1>[[$modelName]]</h1>
//     <button type="button" class="button" hx-get="new" hx-target="#main">New</button>
//     <table>
//         <thead>
//             <tr>
//                 <th>ID</th>[[range .Fields]]
//                 <th>[[.Name]]</th>[[end]]
//                 <th> Actions </th>
//             </tr>
//         </thead>
//         <tbody>{{range .[[.Name]]}}
//             <tr>[[range .Fields]]
//                 <th>{{.[[.Name]]}}</th>[[end]]
//                 <td>
//                     <div class="button-group">
//                         <button type="button" class="button" hx-get="{{.ID}}" hx-target="#main" hx-push-url="true" >Details</button>
//                         <button type="button" class="button warning" hx-get="{{.ID}}/edit" hx-target="#main" hx-push-url="true" >Edit</button>
//                         <button type="button" class="button alert" hx-delete="{{.ID}}" hx-push-url="true">Delete</button>
//                     </div>
//                 </td>
//             </tr>
//         {{end}}</tbody>
//     </table>
// </section>
