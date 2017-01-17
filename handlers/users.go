package handlers

import (
	"html/template"
	"net/http"
	"time"
)

var profile = template.Must(template.New("profile_notifications.html").ParseFiles("templates/profile_notifications.html"))

// ProfileHandler handles the profile page rendering
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: establish session at login or here.

	data := struct {
		Username string
		Time     time.Time
	}{
		Username: "Dummy Name",
		Time:     time.Now(),
	}

	profile.Execute(w, data)
}
