package web

import (
	"log"
	"net/http"

	"time"

	"github.com/btshopng/btshopng/config"
	"github.com/btshopng/btshopng/models"
	uuid "github.com/satori/go.uuid"
)

func NewItemHandler(w http.ResponseWriter, r *http.Request) {

	data, err := GetProfileData(r)
	if err != nil {
		log.Println(err)
		// proper redirect will happen here
		http.Redirect(w, r, "/login", 301)
	}

	tmp := GetTemplates().Lookup("profile_new_barter.html")
	tmp.Execute(w, data)
}

func SaveNewItemHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post data from the request.
	r.ParseForm()
	data, err := GetProfileData(r)
	if err != nil {
		http.Redirect(w, r, "/signup", 301)
		return
	}
	log.Println("Data:", data)

	have := r.FormValue("have")
	haveCat := r.FormValue("haveCat")
	need := r.FormValue("need")
	needCat := r.FormValue("needCat")
	location := r.FormValue("location")

	if have == "" || haveCat == "" || need == "" || needCat == "" || location == "" {
		http.Redirect(w, r, "/newitem?newerror=Fill+out+all+fields", 301)
		return
	}

	uniqueID := uuid.NewV1().String()
	// create a barter model....
	barter := models.Barter{
		ID:           uniqueID,
		UserID:       data.User.ID,
		Have:         have,
		HaveCategory: haveCat,
		Need:         need,
		NeedCategory: needCat,
		Location:     location,
		DateCreated:  time.Now(),
		Status:       "Available",
		Images:       []string{"", "", ""},
	}

	err = barter.Upsert(config.GetConf())
	if err != nil {
		http.Redirect(w, r, "/newitem?error=Could+not+save+your+barter", 301)
		return
	}
	log.Println("New barter added")
}

func ArchiveHandler(w http.ResponseWriter, r *http.Request) {
	data, err := GetProfileData(r)
	if err != nil {
		log.Println(err)
		// proper redirect will happen here
		http.Redirect(w, r, "/login", 301)
	}

	tmp := GetTemplates().Lookup("profile_barter_archive.html")
	tmp.Execute(w, data)
}
