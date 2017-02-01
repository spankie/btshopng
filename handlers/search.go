package handlers

import (
	"html/template"
	"net/http"
)

var listing = template.Must(template.New("list.html").ParseFiles("templates/list.html"))

// SearchHandler lists the search result
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// should place the search result in an array of struct containing a search results
	listing.Execute(w, "")
}
