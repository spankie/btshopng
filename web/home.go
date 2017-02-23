package web

import (
	"log"
	"net/http"

	"github.com/btshopng/btshopng/models"
)

//HomeHandler for listing view
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := Userget(r)
	if err != nil {
		log.Println(err)
	}
	//log.Println(user)

	data := struct {
		User models.User
	}{
		User: user,
	}

	tmp := GetTemplates().Lookup("index.html")
	tmp.Execute(w, data)
}
