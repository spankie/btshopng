package web

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"strings"
)

var templates *template.Template

func TemplateInit() {
	var allFiles []string
	files, err := ioutil.ReadDir("./web/templates")
	if err != nil {
		fmt.Println(err)
	}
	files2, err := ioutil.ReadDir("./web/templates/partials")
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		filename := file.Name()
		log.Println(filename)
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, "./web/templates/"+filename)
		}
	}

	for _, file := range files2 {
		filename := file.Name()
		log.Println(filename)
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, "./web/templates/partials/"+filename)
		}
	}

	templates = template.Must(template.ParseFiles(allFiles...))
}

func GetTemplates() *template.Template {
	return templates
}
