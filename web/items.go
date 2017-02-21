package web

import (
	"encoding/json"
	"log"
	"net/http"

	"time"

	"github.com/btshopng/btshopng/config"
	"github.com/btshopng/btshopng/models"
	uuid "github.com/satori/go.uuid"
)

type BartersListData struct {
	User    models.User
	Barters []models.Barter
}

func NewItemHandler(w http.ResponseWriter, r *http.Request) {

	user, err := Userget(r)
	if err != nil {
		// http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
		log.Printf("User error: %+v", err)
	}
	//log.Println(user)

	result, err := user.Get(config.GetConf())
	if err != nil {
		// http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
		log.Println("\nUser error:", err)
	}

	result.FormattedDateCreated = user.DateCreated.Format("January 2006")

	data := BartersListData{User: result}

	tmp := GetTemplates().Lookup("profile_new_barter.html")
	tmp.Execute(w, data)
}

func SaveNewItemHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post data from the request.
	r.ParseForm()

	user, err := Userget(r)
	if err != nil {
		// http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
		log.Printf("User error: %+v", err)
	}
	conf := config.GetConf()

	result, err := user.Get(conf)
	if err != nil {
		// http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
		log.Printf("User error: %+v", err)
	}

	have := r.FormValue("have")
	haveCat := r.FormValue("haveCat")
	need := r.FormValue("need")
	needCat := r.FormValue("needCat")
	location := r.FormValue("location")
	itemDataImagesString := r.FormValue("itemImageDataInput")

	if have == "" || haveCat == "" || need == "" || needCat == "" || location == "" {
		http.Redirect(w, r, "/newitem?newerror=Fill+out+all+fields", 301)
		return
	}

	itemDataImages := []string{}

	err = json.Unmarshal([]byte(itemDataImagesString), &itemDataImages)
	if err != nil {
		log.Println(err)
	}
	log.Println(itemDataImages)

	Images := []models.Image{}

	for _, bas64Image := range itemDataImages {
		img := models.Image{}
		imgUUID := uuid.NewV1().String()
		fullImageName := "product_images/" + imgUUID
		ThumbnailImageName := "product_thumbnail" + imgUUID
		img.FullSize, err = UploadBase64Image(conf.S3Bucket, bas64Image, fullImageName, 780)
		if err != nil {
			log.Println(err)
		}
		img.Thumbnail, err = UploadBase64Image(conf.S3Bucket, bas64Image, ThumbnailImageName, 340)
		if err != nil {
			log.Println(err)
		}

		Images = append(Images, img)
	}

	uniqueID := uuid.NewV1().String()

	// create a barter model....
	barter := models.Barter{
		ID:           uniqueID,
		UserID:       result.ID,
		Have:         have,
		HaveCategory: haveCat,
		Need:         need,
		NeedCategory: needCat,
		Location:     location,
		DateCreated:  time.Now(),
		Status:       true,
		Images:       Images,
	}

	err = barter.Upsert(conf)
	if err != nil {
		http.Redirect(w, r, "/newitem?error=Could+not+save+your+barter", 301)
		return
	}
	log.Println("New barter added")
	// send a notification to the user that the barter has been added.
	http.Redirect(w, r, "/newitem", 301)
}

func ItemsArchiveHandler(w http.ResponseWriter, r *http.Request) {
	user, err := Userget(r)
	if err != nil {
		// http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
		log.Println(err)
	}

	result, err := user.Get(config.GetConf())
	if err != nil {
		// http.Redirect(w, r, "/signup?loginerror=You+are+not+logged+in", 301)
		log.Println(err)
	}

	// Supply UserID to be used for retrieving all barters.
	barter := models.Barter{UserID: result.ID}

	result.FormattedDateCreated = user.DateCreated.Format("January 2006")
	data := BartersListData{User: result}

	data.Barters, err = barter.GetAll(config.GetConf())
	if err != nil {
		log.Println("No barter for this user.")
	}

	tmp := GetTemplates().Lookup("profile_barter_archive.html")
	tmp.Execute(w, data)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	user, err := Userget(r)
	if err != nil {
		log.Println(err)
	}
	data := struct {
		Barters  []models.Barter
		Query    string
		Location string
		User     models.User
	}{}
	data.User = user

	searchItem := r.URL.Query()

	log.Println("Search item: ", searchItem["searchItem"])
	data.Query = searchItem.Get("searchItem")
	// TODO: GET USER'S GEO LOCATION AS DEFAULT OR USER'S PROFILE LOCATION
	data.Location = "CALABAR"

	data.Barters = []models.Barter{}
	barter := models.Barter{}

	data.Barters, err = barter.GetAllSearch(config.GetConf(), data.Query)
	if err != nil {
		log.Println(err)
	}

	tmp := GetTemplates().Lookup("list.html")
	tmp.Execute(w, data)
}
