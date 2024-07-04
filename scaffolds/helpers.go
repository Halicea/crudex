package scaffolds

import (
	"fmt"
	"os"
	"path/filepath"
)

var cache = map[string]string{}

// ReadContentsOrDefault reads the file if it exists, otherwise returns the default content.
//
// If useCache is true it tries to read the file from the cache first and if it is not
// found there it reads the file from the disk and caches it for any subsequent calls.
func ReadContentsOrDefault(filename, defaultContent string, useCache bool) string {
	fileContent, error := ReadContents(filename, useCache)
	if error == nil {
		return fileContent
	}
	return defaultContent
}

// ReadContents reads the contents of a file and returns it as a string
//
// If useCache is true it tries to read the file from the cache first and if it is not
// found there it reads the file from the disk and caches it for any subsequent calls.
func ReadContents(fname string, useCache bool) (string, error) {
	if useCache {
		if f, ok := cache[fname]; ok {
			return f, nil
		}
	}

	if _, err := os.Stat(fname); err == nil {
		f, err := os.ReadFile(fname)
		if err != nil {
			return "", err
		}
		if useCache {
			cache[fname] = string(f)
			return cache[fname], nil
		} else {
			return string(f), nil
		}
	}
	return "", fmt.Errorf("file not found: %s", fname)
}

func WriteContents(fname, content string, useCache bool) string {
	//check if the directory exists
	if _, err := os.Stat(filepath.Dir(fname)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(fname), 0755)
		if err != nil {
			panic(err)
		}
	}
	if useCache {
		cache[fname] = content
	}
	err := os.WriteFile(fname, []byte(content), 0644)
	if err != nil {
		panic(err)
	}
	return content
}
