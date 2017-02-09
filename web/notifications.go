package web

import (
	"net/http"

	"github.com/btshopng/btshopng/models"
)

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {

	user, err := Userget(r)
	if err != nil {
		http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
	}
	//log.Println(user)

	user.FormattedDateCreated = user.DateCreated.Format("Mon, 02 Jan 2006")

	data := struct {
		User models.User
	}{User: user}

	tmp := GetTemplates().Lookup("profile_notifications.html")
	tmp.Execute(w, data)
}
