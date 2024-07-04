package shared

import (
	"fmt"
)

//go:generate stringer -type=InputKind -linecomment
type InputKind int

const (
	INPUT_TEXT     InputKind = iota //text
	INPUT_MARKDOWN                  //markdown
	INPUT_HTML                      //html
	INPUT_WYSIWYG                   //wysiwyg
	INPUT_TEXTAREA                  //textarea
	INPUT_HIDDEN                    //hidden
	INPUT_PASSWORD                  //password
	INPUT_URL                       //url
	INPUT_EMAIL                     //email
	INPUT_COLOR                     //color
	INPUT_CHECKBOX                  //checkbox
	INPUT_RADIO                     //radio
	INPUT_DATETIME                  //datetime
	INPUT_FILE                      //file
	INPUT_IMAGE                     //image
	INPUT_NUMBER                    //number
	INPUT_RANGE                     //range
	INPUT_SELECT                    //select
	INPUT_SEARCH                    //search
	INPUT_UNKOWN                    //unknown
)

func ParseInputKind(str string) (InputKind, error) {
	prev := 0
	for idx := range _ScaffoldTemplateKind_index {
		if idx == 0 {
			continue
		}
		if str == _ScaffoldTemplateKind_name[prev:_ScaffoldTemplateKind_index[idx]] {
			return InputKind(idx - 1), nil
		}
	}
	return INPUT_UNKOWN, fmt.Errorf("Invalid InputKind string value: %s", str)
}

//go:generate stringer -type=ScaffoldTemplateKind -linecomment
type ScaffoldTemplateKind int

const (
	ScaffoldTemplateLayout  ScaffoldTemplateKind = iota //layout
	ScaffoldTemplateList                                //list
	ScaffoldTemplateDetail                              //detail
	ScaffoldTemplateForm                                //form
	ScaffoldTemplateOpenAPI                             //openapi
)
