package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	//"golang.org/x/oauth2/google"
	"github.com/tonyalaribe/btshopng/config"
	"github.com/tonyalaribe/btshopng/models"
	"golang.org/x/oauth2/facebook"
)

var FBOauthConf = &oauth2.Config{
	ClientID:     "667159983456214",
	ClientSecret: "0a594ec54461df7ecf51406c4d6d44c1",
	Scopes:       []string{"public_profile", "email"},
	RedirectURL:  "http://localhost:8080/fb_oauth_redirect",
	Endpoint:     facebook.Endpoint,
}

// var GoogleOauthConf = &oauth2.Config{
// 	ClientID:     "825438983845-pkg6uce5p4pt7vg74qt7tf8e9850qi2d.apps.googleusercontent.com",
// 	ClientSecret: "0VHYdB6BajL-lRqC_naLOPgV",
// 	Scopes:       []string{"public_profile", "email"},
// 	RedirectURL:  "http://localhost:8080/google_oauth_redirect",
// 	Endpoint:     google.Endpoint,
// }

func SignupPageHandler(w http.ResponseWriter, r *http.Request) {

	data := struct {
		FBAuthURL     string
		GoogleAuthURL string
	}{}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	data.FBAuthURL = FBOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", data)

	// data.GoogleAuthURL := FBOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	// fmt.Printf("Visit the URL for the auth dialog: %v", url)

	tmp := GetTemplates().Lookup("signin_signup.html")
	tmp.Execute(w, data)
}

func FBOauthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if state != "state" {
		http.Redirect(w, r, "/signup", 301)
		return
	}
	ctx := context.Background()
	tok, err := FBOauthConf.Exchange(ctx, code)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signup", 301)
		return
	}
	client := FBOauthConf.Client(ctx, tok)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture,link")
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signup", 301)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signup", 301)
		return
	}
	log.Println(string(b))

	user := models.User{}
	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signup", 301)
		return
	}

	user.Image.URL = user.FBPicture.Data.URL

	user.Upsert(config.GetConf())
	http.Redirect(w, r, "/profile", 301)
}

//
// func GoogleOauthRedirectHandler(w http.ResponseWriter, r *http.Request) {
// 	// Use the authorization code that is pushed to the redirect
// 	// URL. Exchange will do the handshake to retrieve the
// 	// initial access token. The HTTP Client returned by
// 	// conf.Client will refresh the token as necessary.
//
// 	code := r.URL.Query().Get("code")
// 	state := r.URL.Query().Get("state")
// 	if state != "state" {
// 		return
// 	}
// 	ctx := context.Background()
// 	tok, err := GoogleOauthConf.Exchange(ctx, code)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	client := GoogleOauthConf.Client(ctx, tok)
// 	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email")
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	b, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	log.Println(string(b))
// }
