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
	//username := r.FormValue("username")
	//password := r.FormValue("password")

	/*
		db, err := sql.Open("sqlite", db_path)
		sqlstatement := "SELECT id, password_hash, password_salt, email, disabled FROM user WHERE username=" + username
		rows, err := db.Query(sqlstatement)
		if err != nil {
			log.Fatal(err)
		}
		for rows.Next() {
			// Populate User struct
		}
	*/

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
