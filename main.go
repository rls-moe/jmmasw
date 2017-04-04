package main

import (
	"flag"
	"log"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	"os"
	"html/template"
	"bytes"
	"io"
	"strings"
)

var (
	projectDir = flag.String("dir", ".", "Project directory")
	dataFile = flag.String("data", "./data.json", "Variable Data File")
	failMissingData = flag.Bool("ignore-no-data", true, "Ignore a missing Data Files")
	outDir = flag.String("static-out", "./jasw-out", "Static Output of Website")
)

func main() {
	flag.Parse()
	t, err := LoadTemplatesAtPath(*projectDir)
	if err != nil {
		log.Fatal("Error in template loader: ", err)
		return
	}
	var data interface{}
	if !*failMissingData {
		if _,err := os.Stat(*dataFile); os.IsNotExist(err) {
			log.Fatal("Data file missing")
			return
		}
	}
	dat, err := ioutil.ReadFile(*dataFile)
	if err != nil {
		log.Println("No data file defined or error opening, skipping")
		log.Println("Error: ", err)
		data = nil
	} else {
		err = json.Unmarshal(dat, &data)
		if err != nil {
			log.Fatal("Error on parsing datafile: ", err)
			return
		}
	}
	outDirAbs, err := filepath.Abs(*outDir)
	if _, err := os.Stat(outDirAbs); os.IsNotExist(err) {
		err = os.MkdirAll(outDirAbs, 0700)
		if err != nil {
			log.Fatal("Error while creating output directory: ", err)
			return
		}
	}

	log.Println("Running all Templates")
	for _, v := range t.Templates() {
		func (v *template.Template) {
			if v.Name() == "root" {
				return
			}
			var buf = bytes.NewBuffer([]byte{})
			err = v.Execute(buf, data)
			if err != nil {
				log.Fatal("Error executing template: ", err)
				return
			}

			if len(strings.TrimSpace(buf.String())) == 0 {
				log.Fatal("Template ", v.Name(), " returned empty file, skipping output")
				return
			}

			curFile := filepath.Join(outDirAbs, v.Name())
			if _, err := os.Stat(filepath.Dir(curFile)); os.IsNotExist(err) {
				if err := os.MkdirAll(filepath.Dir(curFile), 0700); err != nil {
					log.Fatal("Error while creating dir structure: ", err)
					return
				}
			}
			file, err := os.Create(curFile)
			if err != nil {
				log.Fatal("Error creating file: ", err)
				return
			}
			defer file.Close()

			_, err = io.Copy(file, buf)
			if err != nil {
				log.Fatal(err)
			}
		}(v)
	}
	log.Println("Fin!")
}