package web

import (
	"log"
	"net/http"

	"github.com/btshopng/btshopng/models"
)

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {

	user, err := Userget(r)
	if err != nil {
		// http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
		log.Println("User error:", err)
	}
	//log.Println(user)

	user.FormattedDateCreated = user.DateCreated.Format("January 2006")

	data := struct {
		User models.User
	}{User: user}
	//log.Printf("%+v", data)
	tmp := GetTemplates().Lookup("profile_notifications.html")
	tmp.Execute(w, data)
}
