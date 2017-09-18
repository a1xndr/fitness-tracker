package db

import (
	"github.com/jmoiron/sqlx"
	"log"
)

var SQL *sqlx.DB

type DBinfo struct {
	Path string
}

func Connect(d DBinfo) {
	var err error
	SQL, err = sqlx.Open("sqlite3", d.Path)
	if err != nil {
		log.Fatal(err)
	}
}
