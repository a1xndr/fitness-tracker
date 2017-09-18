package models

import (
	"alpha/db"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	Id             uint32    `db:"id"`
	Username       string    `db:"username"`
	PasswordHashed string    `db:"password_hashed"` //bcrypt
	Email          string    `db:"email"`
	CreatedAt      time.Time `db:"created_at"`
	Disabled       bool      `db:"disabled"`
}

var db_path string = "./alpha.db"

func UserCreate(username string, email string, password string) {
	var err error
	time := time.Now()
	db.SQL.Exec(
		`INSERT INTO user(username,password_hashed,email,created_at,disabled)
		VALUES (?,?,?,?,?)
		`,
		username,
		HashAndSalt(password),
		email,
		time,
		false,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func UserByUsername(username string) (User, error) {
	var err error

	user := User{}

	err = db.SQL.Get(&user,
		`SELECT id,username,password_hashed,email,created_at,disabled 
		FROM user
		WHERE username = ? 
		LIMIT 1
	`, username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func HashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func (u *User) CheckPasswordMatch(password string) bool {
	result := bcrypt.CompareHashAndPassword([]byte(u.PasswordHashed),
		[]byte(password))

	return result == nil
}
