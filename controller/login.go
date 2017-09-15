package controller

import (
	"html/template"
	"log"
	"net/http"
)

func LoginGET(w http.ResponseWriter, r *http.Request) {
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
//  username := r.FormValue("username")
//	password := r.FormValue("password")
	

	// Check the password against hash + salt from DB

	// Redirect them to the dashboard or reload the page
	tmpl := template.Must(template.ParseFiles(
		"templates/base/layout.html",
		"templates/login.html",
	))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}

}
