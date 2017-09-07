package model

import "fmt"

type User struct {
	Id            uint32 `db:"id"`
	username      string `db:"username"`
	password_hash string `db:"password_hash"`
	password_salt string `db:"password_salt"`
	email         string `db:"email"`
	disabled      bool   `db:"disabled"`
}
