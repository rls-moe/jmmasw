package main

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	extWhitelist = []string{
		".html", ".tmpl",
	}
)

var (
	functions = map[string]interface{}{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"markdown": func(text string) template.HTML {
			return template.HTML(blackfriday.MarkdownCommon([]byte(text)))
		},
		"json": func(jsonStr string) (map[string]interface{}, error) {
			var dat = map[string]interface{}{}
			err := json.Unmarshal([]byte(jsonStr), &dat)
			return dat, err
		},
		"file": func(text string) (string, error) {
			dat, err := ioutil.ReadFile(text)
			return string(dat), err
		},
	}
)

func LoadTemplatesAtPath(path string) (*template.Template, error) {
	var mainTemplate = template.New("root")

	mainTemplate.Funcs(functions)
	mainTemplate.Option("missingkey=error")

	path = filepath.Clean(path)

	if stat, err := os.Stat(path); os.IsNotExist(err) {
		return mainTemplate, err
	} else if !stat.IsDir() {
		return mainTemplate, errors.New("template path must be a directory")
	}

	outDirRel, err := filepath.Rel(path, filepath.Clean(*outDir))
	if err != nil {
		return nil, err
	}

	filepath.Walk(path, func(fpath string, stat os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !func(ext string) bool {
			for _, v := range extWhitelist {
				if v == ext {
					return true
				}
			}
			return false
		}(filepath.Ext(fpath)) {
			log.Print("Skipping ", fpath, " since it's not a templateable file")
			return nil
		}

		tmplName, err := filepath.Rel(path, fpath)
		if err != nil {
			return err
		}

		if strings.HasPrefix(tmplName, outDirRel+"/") {
			log.Print("Skipping ", tmplName, " since it's in the output folder")
			return nil
		}

		if !stat.IsDir() {
			dat, err := ioutil.ReadFile(fpath)
			if err != nil {
				return err
			}
			_, err = mainTemplate.New(tmplName).Parse(string(dat))
			return err
		}
		return nil
	})

	for _, v := range mainTemplate.Templates() {
		log.Println("Read Template: ", v.Name())
	}

	return mainTemplate, nil
}
