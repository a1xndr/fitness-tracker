package controller

import (
	"alpha/models"
	"html/template"
	"log"
	"net/http"
)

func LoginGET(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		LoginPOST(w, r)
		return
	}
	// View
	tmpl := template.Must(template.ParseFiles(
		"templates/base/layout.html",
		"templates/login.html",
	))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := models.UserByUsername(username)

	if err != nil {
		http.Redirect(w, r, "/login/", 307)
		log.Fatal(err)
	}
	if user.CheckPasswordMatch(password) {
		http.Redirect(w, r, "/dashboard/", 307)
	}
}
