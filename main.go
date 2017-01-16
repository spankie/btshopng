package main

import (
	"log"
	"net/http"
	"os"

	"github.com/spankie/btshopng/config"
	"github.com/spankie/btshopng/handlers"
)

var (
	appConf *config.Conf
)

func main() {

	// Initialize configurations
	config.Init()
	// Close the session
	defer config.GetConf().Database.Session.Close()

	// Serve the signup page.
	http.HandleFunc("/signup", handlers.SignupHandler)
	// process login form
	http.HandleFunc("/login", handlers.LoginHandler)
	// Serve profile page
	http.HandleFunc("/profile", handlers.ProfileHandler)

	// Serve static files.
	f := http.FileServer(http.Dir("./templates/assets/"))
	http.Handle("/public/", http.StripPrefix("/public/", f))

	// Start up the server
	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Println("No port set, Using default port 3000")
		PORT = "3000"
	}

	log.Printf("Server started on %s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
