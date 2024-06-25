package crudex

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

type ScaffoldStrategy int

func RenderTypeInput(modelName string, field reflect.StructField) string {
	switch field.Type.Kind() {
	case reflect.String:
		return fmt.Sprintf(`<input type="text" name="%s">{{.%s.%s}}</input>`, field.Name, modelName, field.Name)
	case reflect.Int, reflect.Float64:
		return fmt.Sprintf(`<input type="number" name="%s">{{.%s.%s}}</input>`, field.Name, modelName, field.Name)
	case reflect.Bool:
		return fmt.Sprintf(`<input type="checkbox" name="%s">{{.%s.%s}}</input>`, field.Name, modelName, field.Name)
	}
	panic(fmt.Sprintf("unsupported type: %s for field %s", field.Type.Kind().String(), field.Name))
}

var cache = map[string]string{}

func readContents(fname string) (string, error) {
	if f, ok := cache[fname]; ok {
		return f, nil
	}
	if _, err := os.Stat(fname); err == nil {
		f, err := os.ReadFile(fname)
		if err != nil {
			return "", err
		}
		cache[fname] = string(f)
		return cache[fname], nil
	}
	return "", fmt.Errorf("file not found: %s", fname)
}

func WriteContents(fname, content string) string {
	//check if the directory exists
	if _, err := os.Stat(filepath.Dir(fname)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(fname), 0755)
		if err != nil {
			panic(err)
		}
	}

	cache[fname] = content
	err := os.WriteFile(fname, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
	return content
}
