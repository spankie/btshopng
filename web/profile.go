package web

import (
	"log"
	"net/http"

	"github.com/btshopng/btshopng/config"
	"github.com/btshopng/btshopng/models"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user, err := Userget(r)
	if err != nil {
		log.Println(err)
	}
	//log.Println(user)

	data := struct {
		User models.User
	}{}
	data.User, err = user.Get(config.GetConf())
	if err != nil {
		log.Println(err)
	}

	data.User.FormattedDateCreated = user.DateCreated.Format("January 2006")
	tmp := GetTemplates().Lookup("profile_new_barter.html")
	tmp.Execute(w, data)
}
