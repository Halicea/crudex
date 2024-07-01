package scaffolds

import (
	"reflect"
	"strings"
	"testing"
)

type TestStruct struct {
	Str  string
	Num  int32
	Html string `crud-input:"html" crud-placeholder:"Enter some HTML"`
	Date string `crud-input:"datetime"`
}

func TestRender_InputTypes(t *testing.T) {
	var exprectedTypes = map[string]string{
		"Str":  "type=\"text\"",
		"Num":  "type=\"number\"",
		"Html": "type=\"textarea\"",
		"Date": "type=\"datetime\"",
	}
	tt := reflect.TypeFor[TestStruct]()
	for i := 0; i < tt.NumField(); i++ {
		field := tt.FieldByIndex([]int{i})
		res := RenderInputType(tt.Name(), field)
		t.Logf("Output: %s", res)
		if !strings.Contains(res, exprectedTypes[field.Name]) {
			t.Errorf("Expected %s, got %s", exprectedTypes[field.Name], res)
		}
	}
}
