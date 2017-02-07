package web

import (
	"log"
	"net/http"
)

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := GetProfileData(r)
	if err != nil {
		log.Println(err)
		// proper redirect will happen here
		http.Redirect(w, r, "/login", 301)
	}

	tmp := GetTemplates().Lookup("profile_notifications.html")
	tmp.Execute(w, data)
}
