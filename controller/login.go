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
