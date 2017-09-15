package controller

import (
	"html/template"
	"log"
	"net/http"
	"alpha/models"
)

func RegisterGET(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login/", 307)
}

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	models.UserCreate(username, email, password)

	tmpl := template.Must(template.ParseFiles(
		"templates/base/layout.html",
		"templates/register.html",
	))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}

}
