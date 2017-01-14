package main

import (
	"fmt"
	"html/template"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// template of signup/signin page to be served
var tmpl = template.Must(template.New("signin_signup.html").ParseFiles("../signin_signup.html"))

// User interface to contains user information
type User struct {
	name     string
	email    string
	password string
}

func main() {

	// Serve the signup page.
	http.HandleFunc("/signup", signupHandler)

	// Serve static files.
	f := http.FileServer(http.Dir("../"))
	http.Handle("/public/", http.StripPrefix("/public/", f))

	// Start up the server
	port := "3000"
	fmt.Printf("Server started on %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {

	// check if the request is a post request.
	if r.Method == "POST" {
		// If it is a post request, process the request.

		// connect to mongodb server
		var session, err = mgo.Dial("mongodb://spankie:506dad@ds163738.mlab.com:63738/btshopng")
		if err != nil {
			// reply with internal server error
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// get the post form values
		r.ParseForm()
		// new struct containing the form values
		newUser := &User{name: r.PostFormValue("name"), email: r.PostFormValue("email"), password: r.PostFormValue("passwd")}

		// create a db connection
		c := session.DB("btshopng").C("Users")

		// check if email has already been used by querying the db
		result := User{}
		err = c.Find(bson.M{"email": newUser.email}).Select(bson.M{"email": 0}).One(&result)
		if err == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// compare the result from the db to an empty struct
		if result == (User{}) {

			// if result is empty, then the email can be used.
			fmt.Println("Email is available")

			// insert the User data to db
			err = c.Insert(newUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				data := "inserted"
				tmpl.Execute(w, data)
			}

		} else {

			// if result is not empty, then the email has already been used before.
			fmt.Println("email already exists")
			data := "email already exists"
			tmpl.Execute(w, data)
		}

	} else if r.Method == "GET" {
		// if the request is not a post request, just Serve the page
		data := "name"
		tmpl.Execute(w, data)
	}

}
