package web

import (
	"log"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user, err := Userget(r)
	if err != nil {
		log.Println(err)
	}
	tmp := GetTemplates().Lookup("profile_new_barter.html")
	tmp.Execute(w, "")
}
