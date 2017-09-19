package models

import (
	"alpha/db"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type Exercise struct {
	Id          uint32 `db:"id"`
	Name        string `db:"Name"`
	TypeId      uint32 `db:"TypeId"` //bcrypt
	Description string `db:"description"`
	Reps        bool   `db:"reps"`
	Weight      bool   `db:"weight"`
	Time        bool   `db:"time"`
	Speed       bool   `db:"speed"`
	Grade       bool   `db:"grade"`
}
