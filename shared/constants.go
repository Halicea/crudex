package shared

const (
	INPUT_TEXT     = "text"
	INPUT_MARKDOWN = "markdown"
	INPUT_HTML     = "html"
	INPUT_WYSIWYG  = "wysiwyg"
	INPUT_TEXTAREA = "textarea"
	INPUT_HIDDEN   = "hidden"
	INPUT_PASSWORD = "password"
	INPUT_URL      = "url"

	INPUT_EMAIL = "email"
	INPUT_COLOR = "color"

	INPUT_CHECKBOX = "checkbox"
	INPUT_RADIO    = "radio"

	INPUT_DATETIME = "datetime"

	INPUT_FILE  = "file"
	INPUT_IMAGE = "image"

	INPUT_NUMBER = "number"
	INPUT_RANGE  = "range"

	// for relations searching
	INPUT_SELECT = "select"
	INPUT_SEARCH = "search"
)


const (
    ScaffoldTemplateLayout = "layout"
    ScaffoldTemplateList = "list"
    ScaffoldTemplateDetail = "detail"
    ScaffoldTemplateForm = "form"
)
var AllScaffoldTemplates = []string{ScaffoldTemplateLayout, ScaffoldTemplateList, ScaffoldTemplateDetail, ScaffoldTemplateForm}
