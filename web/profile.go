package web

import (
	"log"
	"net/http"

	"github.com/tonyalaribe/btshopng/config"
	"github.com/tonyalaribe/btshopng/models"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user, err := Userget(r)
	if err != nil {
		log.Println(err)
	}
	log.Println(user)

	data := struct {
		User models.User
	}{}
	data.User, err = user.Get(config.GetConf())
	if err != nil {
		log.Println(err)
	}
	log.Println(data)
	tmp := GetTemplates().Lookup("profile_new_barter.html")
	tmp.Execute(w, data)
}
