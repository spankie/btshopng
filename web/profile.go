package web

import "net/http"

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	tmp := GetTemplates().Lookup("profile_new_barter.html")
	tmp.Execute(w, "")
}
