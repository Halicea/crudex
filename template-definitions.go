package crudex
const tmplLayout = `{{/* generated file: [[.TemplateFileName]] */}}
<!doctype html>
<html>
  <head>
    <title>Index</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta charset="UTF-8">

    <script src="https://unpkg.com/htmx.org@1.7.0" integrity="sha384-EzBXYPt0/T6gxNp0nuPtLkmRpmDBbjg6WmCUZRLXBBwYYmwAUxzlSGej0ARHX0Bo" crossorigin="anonymous" defer></script>
    {{template "style.html" .}}
  </head>
  <body>
    <header>
      <nav>
        <ul class="menu">
          <li><a href="{{.Path}}">Home</a></li>
          [[range .Menu]]<li><a href="#" hx-get="[[ .Path ]]" hx-target="#main" hx-push-url="true">[[ .Title ]]</a></li>
          [[end]]
        </ul>
      </nav>
    </header>
    <main> 
        <div class="content" id="main">
        </div>
    </main>
    <footer>
      <p>Footer</p>
    </footer>
  </body>
</html>`

const tmplDetail = `{{/* generated file: [[.TemplateFileName]] */}}
<section>
    [[$modelName := .Name]]
    <h1>[[$modelName]]</h1>
    <div>[[range .Fields]]
        <div>
            <label for="[[.Name]]">[[.Name]]</label>
            <div>{{.[[$modelName]].[[.Name]]}}</div>
        </div>
    [[end]]</div>
</section>`

const tmplList = `{{/* generated file: [[.TemplateFileName]] */}}
<section>
    [[$modelName := .Name]]
    <h1>[[$modelName]]</h1>
    <button hx-get="new" hx-target="#main">New</button>
    <table>
        <thead>
            <tr>[[range .Fields]]
                <th>[[.Name]]</th>[[end]]
                <th> Actions </th>
            </tr>
        </thead>
        <tbody>{{range .[[.Name]]}}
            <tr>[[range .Fields]]
                <th>{{.[[.Name]]}}</th>[[end]]
                <td>
                    <button hx-get="{{.[[$modelName]].ID}}" hx-target="#main">Details</button>
                    <button hx-get="{{.[[$modelName]].ID}}/edit" hx-target="#main">Edit</button>
                    <button hx-delete="{{.[[$modelName]].ID}}">Delete</button>
                </td>
            </tr>
        {{end}}</tbody>
    </table>
</section>`

const tmplForm = `{{/* generated file: [[.TemplateFileName]] */}}
<section>
    [[$modelName := .Name]]
    <h1>[[$modelName]]</h1>
    <form
        {{if .ID}}hx-post="{{.Path}}"{{else}}hx-put="{{.Path}}"{{end}}
        hx-target="#main">[[range .Fields]]
        <div>
            <label for="[[.Name]]">[[.Name]]</label>
            [[RenderTypeInput .]]
        </div>[[end]]
        <button type="submit">Submit</button>
    </form>
<section>`
