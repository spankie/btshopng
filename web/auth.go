package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	//"golang.org/x/oauth2/google"
	"github.com/btshopng/btshopng/config"
	"github.com/btshopng/btshopng/models"
)

//LoginResponse sent to the cllient, carrying the token, upon login
type LoginResponse struct {
	User    models.User
	Message string
	Token   string
}

// spankie's oauthConfig
// var FBOauthConf = &oauth2.Config{
// 	ClientID:     "763167067164923",
// 	ClientSecret: "9fadba8f65774f03d492ca95128e1a09",
// 	Scopes:       []string{"public_profile", "email"},
// 	RedirectURL:  "http://localhost:8080/fb_oauth_redirect",
// 	Endpoint:     facebook.Endpoint,

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
		LoginMessage  string
		SignupError   string
	}{}

	data.LoginMessage = ""
	data.SignupError = ""

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	data.FBAuthURL = FBOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", data)

	// data.GoogleAuthURL := FBOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	// fmt.Printf("Visit the URL for the auth dialog: %v", url)

	tmp := GetTemplates().Lookup("signin_signup.html")
	tmp.Execute(w, data)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// check if POST data are set and validate them
	// I feel SignupPageHandler() and this function should share the same struct for better memory conservation
	data := struct {
		FBAuthURL     string
		GoogleAuthURL string
		LoginMessage  string
		SignupError   string
	}{}

	data.LoginMessage = ""
	data.SignupError = ""

	// This is set here so that when there are any errors from the signup process,
	// the link will be passed to the template alongside the errors.
	data.FBAuthURL = FBOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	r.ParseForm()

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" && password == "" {

		data.LoginMessage = "Please Fill out all required Fields"
		tmp := GetTemplates().Lookup("signin_signup.html")
		tmp.Execute(w, data)
		return

	}

	user := models.User{
		Email: email,
	}

	// Thinking of making this a standalone function (CheckUser()) for proper memory management
	result, err := user.CheckUser(config.GetConf())

	if err != nil {
		// data.LoginMessage = err.Error()
		data.LoginMessage = "Username or password incorrect"
		tmp := GetTemplates().Lookup("signin_signup.html")
		tmp.Execute(w, data)
		return
	}

	err = bcrypt.CompareHashAndPassword(result.Password, []byte(password))
	if err != nil {
		// password is incorrect
		data.LoginMessage = "Username or password incorrect"
		log.Println("db pass: ", result.Password, "form pass: ", password)
		// log.Println("User: ", result)
		tmp := GetTemplates().Lookup("signin_signup.html")
		tmp.Execute(w, data)
		return
	}

	// data.LoginMessage = "You should be logged in by now"
	// tmp := GetTemplates().Lookup("signin_signup.html")
	// tmp.Execute(w, data)
	// log.Println("User: ", result)

	//// THE NEXT FEW LINES PROBABLY SHOULD HAPPEN IN ANOTHER FUNCTION SINCE IT IS USED BY THESE THREE HANDLERS
	// LOGINHANDLER, SIGNUP HANDLER AND FB HANDLER.

	// If all goes well, generate a token for the cookie
	loginResp, err := GenerateJWT(result)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signup", 301)
		return
	}
	// Create cookie
	expire := time.Now().AddDate(0, 0, 1)
	// I don't really know if the Name of the token should change
	// from the one you used at FBOauthRedirectHandler()
	cookie := http.Cookie{Name: "AuthToken", Value: loginResp.Token, Path: "/", Expires: expire, MaxAge: 86400}
	// Set cookie
	http.SetCookie(w, &cookie)
	// send the user to thier profile page
	http.Redirect(w, r, "/profile", 301)

	// **** user.Get() returns the same user unchanged ****

	// Check if email matches any in the DB
	// nuser, err := user.Get(config.GetConf())
	// if err != nil {
	// 	log.Println("LOGIN ERROR: ", err, " USER: ", nuser)
	// 	data.LoginMessage = "Could not Sign you in"
	// 	tmp := GetTemplates().Lookup("signin_signup.html")
	// 	tmp.Execute(w, data)
	// 	return
	// }

	// log.Println("USER: ", nuser)
	// data.LoginMessage = "You are logged in"
	// tmp := GetTemplates().Lookup("signin_signup.html")
	// tmp.Execute(w, data)
}

// SignupHandler handles the signup process for in app signup
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// I feel SignupPageHandler() and this function should share the same struct for better memory conservation
	data := struct {
		FBAuthURL     string
		GoogleAuthURL string
		LoginMessage  string
		SignupError   string
	}{}

	data.LoginMessage = "Login"
	data.SignupError = ""

	// This is set here so that when there are any errors from the signup process,
	// the link will be passed to the template alongside the errors.
	data.FBAuthURL = FBOauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	r.ParseForm()

	fullName := r.FormValue("name")
	email := r.FormValue("email")
	passwd := r.FormValue("passwd")

	// Check if the Post data not empty and validate them.
	if fullName == "" && email == "" && passwd == "" {

		data.SignupError = "Please Fill out all required Fields"
		tmp := GetTemplates().Lookup("signin_signup.html")
		tmp.Execute(w, data)
		return

	}

	// encrypt password.
	password, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	now := time.Now()

	// create the user data
	user := models.User{
		Name:                 fullName,
		Email:                email,
		DateCreated:          now,
		FormattedDateCreated: now.String(),
		Password:             password,
	}

	// CHECK IF THE EMAIL HAS ALREADY BEEN USED BEFORE.
	_, err := user.CheckUser(config.GetConf())
	if err != nil {
		data.SignupError = "Email address has already been used."
		tmp := GetTemplates().Lookup("signin_signup.html")
		tmp.Execute(w, data)
		return
	}
	// *** Upsert is replacing documents on the collection. ***

	// Upsert the user data to the db
	err = user.Upsert(config.GetConf())
	if err != nil {
		log.Println(err)
		data.SignupError = "Could not Sign you up right now. Try Again"
		tmp := GetTemplates().Lookup("signin_signup.html")
		tmp.Execute(w, data)
		return
	}

	// note: Userget(r) is passing a string to User.Password instead of []byte

	// SHOULD VERIFY EMAIL ADDRESS SENT, HERE.

	// If all goes well, generate a token for the cookie
	loginResp, err := GenerateJWT(user)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signup", 301)
		return
	}
	// Create cookie
	expire := time.Now().AddDate(0, 0, 1)
	// I don't really know if the Name of the token should change
	// from the one you used at FBOauthRedirectHandler()
	cookie := http.Cookie{Name: "AuthToken", Value: loginResp.Token, Path: "/", Expires: expire, MaxAge: 86400}
	// Set cookie
	http.SetCookie(w, &cookie)
	// send the user to thier profile page
	http.Redirect(w, r, "/profile", 301)
	// return

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
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture.type(large),link")
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
	user.DateCreated = time.Now()
	err = user.Upsert(config.GetConf())
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signup", 301)
		return
	}

	loginResp, err := GenerateJWT(user)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signup", 301)
		return
	}

	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{Name: "AuthToken", Value: loginResp.Token, Path: "/", Expires: expire, MaxAge: 86400}

	http.SetCookie(w, &cookie)

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
