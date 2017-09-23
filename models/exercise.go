package models

import (
        "alpha/db"
        "log"
)

type Exercise struct {
        Id          uint32 `db:"id"`
        Name        string `db:"name"`
        Type        uint32 `db:"type"`
        Description string `db:"description"`
        Reps        bool   `db:"reps"`
        Weight      bool   `db:"weight"`
        Time        bool   `db:"time"`
        Speed       bool   `db:"speed"`
        Grade       bool   `db:"grade"`
}

func (e *Exercise) Store() {
        var err error
        db.SQL.Exec(
                `INSERT OR REPLACE INTO exercise
        (id, name, type, description, reps, weight, time, speed, grade
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
                e.Id,
                e.Name,
                e.Type,
                e.Description,
                e.Reps,
                e.Weight,
                e.Time,
                e.Speed,
                e.Grade,
        )
        if err != nil {
                log.Fatal(err)
        }

}

