package handlers

import (
	"html/template"
	"log"
	"net/http"

	"time"

	"github.com/spankie/btshopng/config"
	"gopkg.in/mgo.v2/bson"
)

// template of signup/signin page to be served
var tmpl = template.Must(template.New("signin_signup.html").ParseFiles("templates/signin_signup.html"))

// Data contains data to be passed to templates
type Data struct {
	SignupMessage string
	LoginMessage  string
}

// User interface to contain user information
type User struct {
	Name     string
	Email    string
	Password string
	Time     time.Time
}

// LoginHandler handles Login process
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// instantiate data
	data := Data{}

	if r.Method == "POST" {

		// Get database from configurations
		db := config.GetConf().Database
		// db := appConf.Database

		// get the form values
		r.ParseForm()

		// TODO: Validate form data before checking the db.(ifempty, matches the required format)
		email := r.PostFormValue("email")
		// TODO: password would be hashed to match the DB result
		passwd := r.PostFormValue("password")

		// select Collection
		c := db.C("Users")

		// result struct
		result := bson.D{}

		// Check if email and password matches any in the DB
		err := c.Find(bson.M{"Email": email, "Password": passwd}).One(&result)
		if err != nil {

			// debug purposes
			log.Println("result:", result)

			w.Header().Set("Content-Type", "text/html")
			data.LoginMessage = "Username or Password is incorrect"
			tmpl.Execute(w, data)

			// http.Error(w, err.Error(), http.StatusInternalServerError)
			// log.Println("no data matching the mail supplied", err)

		} else {
			// TODO: redirect to the the users profile page with session
			http.Redirect(w, r, "/profile", http.StatusFound)
		}

	} else if r.Method == "GET" {

		// if the request is not a post request, just Serve the page
		tmpl.Execute(w, data)

	}

}

// SignupHandler handles signup process
func SignupHandler(w http.ResponseWriter, r *http.Request) {

	// instantiate data
	data := Data{}

	// check if the request is a post request.
	if r.Method == "POST" {
		// If it is a post request, process the request.

		db := config.GetConf().Database

		// get the post form values
		r.ParseForm()
		// new struct containing the form values
		newUser := User{
			Name:  r.PostFormValue("name"),
			Email: r.PostFormValue("email"),
			// TODO: Password would be hashed before storage
			Password: r.PostFormValue("passwd"),
			Time:     time.Now(),
		}

		// create a db connection
		c := db.C("Users")

		// check if email has already been used by querying the db
		var count int
		count, err := c.Find(bson.M{"Email": string(newUser.Email)}).Count()
		// Select(bson.M{"email": 0}).
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// compare the result from the db to an empty struct
		if count <= 0 {

			// if result is empty, then the email can be used.
			log.Println("Email is available")

			// insert the User data to db
			err = c.Insert(newUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			// Execute template
			data.SignupMessage = "A verification email has been sent to you"
			tmpl.Execute(w, data)

		} else {

			// if result is not empty, then the email has already been used before.
			log.Println("email already exists")
			data.SignupMessage = "Email already exists"
			tmpl.Execute(w, data)
		}

	} else if r.Method == "GET" {
		// if the request is not a post request, just Serve the page
		tmpl.Execute(w, data)
	}

}
