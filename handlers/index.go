package handlers

import (
	"html/template"
	"net/http"
)

var index = template.Must(template.New("index.html").ParseFiles("templates/index.html"))

// IndexHandler renders the index page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	index.Execute(w, "")
}
