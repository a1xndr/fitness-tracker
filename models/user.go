package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"time"
	"log"
)

type User struct {
	Id             uint32
	Username       string
	PasswordHashed string // And salted
	Email          string
	CreatedAt      time.Time
	Disabled       bool
}

var db_path string = "./alpha.db"

func UserCreate(username string, email string, password string) {
	var err error
	time := time.Now()

	db, err := sql.Open("sqlite3", db_path)
	sqlstatement, err := db.Prepare(
		`INSERT INTO user(username,password_hashed,email,created_at,disabled)
		VALUES (?,?,?,?,?,?)
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = sqlstatement.Exec(
		username,
		hashAndSalt(password),
		email,
		time,
		false,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func hashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}
