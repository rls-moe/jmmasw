package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	projectDir      = flag.String("dir", ".", "Project directory")
	dataFile        = flag.String("data", "./data.json", "Variable Data File")
	failMissingData = flag.Bool("ignore-no-data", true, "Ignore a missing Data Files")
	outDir          = flag.String("static-out", "./jasw-out", "Static Output of Website")
)

const (
	folderPerm = 0755
	filePerm   = 0644
)

func main() {
	flag.Parse()
	t, err := LoadTemplatesAtPath(*projectDir)
	if err != nil {
		log.Fatal("Error in template loader: ", err)
		return
	}
	var data = map[string]interface{}{}
	if !*failMissingData {
		if _, err := os.Stat(*dataFile); os.IsNotExist(err) {
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
		err = os.MkdirAll(outDirAbs, folderPerm)
		if err != nil {
			log.Fatal("Error while creating output directory: ", err)
			return
		}
	}

	log.Println("Running all Templates")
	for _, v := range t.Templates() {
		func(v *template.Template) {
			if v.Name() == "root" {
				return
			}
			if strings.HasSuffix(v.Name(), ".tmpl") {
				return
			}
			var buf = bytes.NewBuffer([]byte{})
			err = v.Execute(buf, data)
			if err != nil {
				log.Print("Error executing template: ", err)
				if err, ok := err.(*template.Error); ok {
					log.Fatal("Error on line ", err.Line,
						"\n", err.Description)
				}
				return
			}

			if len(strings.TrimSpace(buf.String())) == 0 {
				log.Fatal("Template ", v.Name(), " returned empty file, skipping output")
				return
			}

			curFile := filepath.Join(outDirAbs, v.Name())
			if _, err := os.Stat(filepath.Dir(curFile)); os.IsNotExist(err) {
				if err := os.MkdirAll(filepath.Dir(curFile), 0755); err != nil {
					log.Fatal("Error while creating dir structure: ", err)
					return
				}
			}
			//file, err := os.Create(curFile)
			file, err := os.OpenFile(curFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePerm)
			if err != nil {
				log.Fatal("Error creating file: ", err)
				return
			}
			defer file.Close()

			_, err = io.Copy(file, buf)
			if err != nil {
				log.Fatal(err)
				return
			}
		}(v)
	}
	log.Println("Fin!")
}
