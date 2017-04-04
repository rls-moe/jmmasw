package main

import (
	"html/template"
	"path/filepath"
	"os"
	"github.com/go-errors/errors"
	"io/ioutil"
	"log"
	"strings"
)

var (
	extWhitelist = []string{
		".html", ".js", ".css", ".tmpl",
	}
)

func LoadTemplatesAtPath(path string) (*template.Template, error) {
	var mainTemplate = template.New("root")

	path = filepath.Clean(path)

	if stat, err := os.Stat(path); os.IsNotExist(err) {
		return mainTemplate, err
	} else if !stat.IsDir() {
		return mainTemplate, errors.New("template path must be a directory")
	}

	outDirRel,err := filepath.Rel(path, filepath.Clean(*outDir))
	if err != nil {
		return nil, err
	}

	filepath.Walk(path, func (fpath string, stat os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !func(ext string) bool {
			for _,v := range extWhitelist {
				if v == ext {
					return true
				}
			}
			return false
		}(filepath.Ext(fpath)) {
			return nil
		}

		tmplName, err := filepath.Rel(path, fpath)
		if err != nil {
			return err
		}

		if strings.HasPrefix(tmplName, outDirRel + "/") {
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