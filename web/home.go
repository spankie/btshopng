package web

import "net/http"

//HomeHandler for listing view
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmp := GetTemplates().Lookup("index.html")
	tmp.Execute(w, "")
}
