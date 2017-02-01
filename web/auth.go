package web

import "net/http"

func SignupPageHandler(w http.ResponseWriter, r *http.Request) {
	tmp := GetTemplates().Lookup("signin_signup.html")
	tmp.Execute(w, "")
}
