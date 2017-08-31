package controllers

import (
	"alpha/web"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ServiceController struct {
	BaseURL  string
	Renderer func(http.ResponseWriter, *http.Request)
}

func (sc *ServiceController) Register(s *web.Server) {
	s.HandleFunc(sc.BaseURL, sc.Renderer)
}

/* O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_O_o_ */
var db_path string = "./alpha.db"

type Set struct {
	Id       uint64
	Exercise string
	Reps     uint64
	Weight   float64
	Sets     uint64
	Seconds  float64
}
type Workout struct {
	Time time.Time
	Id   uint64
	Sets []Set
}

type Exercise struct {
	Id          uint64
	Name        string
	Description string
	Reps        bool
	Weight      bool
	Seconds     bool
	Speed       bool
	Grade       bool
}
type Context struct {
	Workout   *Workout
	Exercises []Exercise
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
func (w *Workout) SaveWorkout() error {
	return w.SaveSetInDB()
}

func (ex *Exercise) SaveExercise() error {
	db, err := sql.Open("sqlite3", db_path)
	sqlstatement, err := db.Prepare(`
            INSERT INTO exercise(name, description, reps, weight, seconds, speed, grade)
            VALUES(?,?,?,?,?,?,?)
            `)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = sqlstatement.Exec(
		ex.Name,
		ex.Description,
		ex.Reps,
		ex.Weight,
		ex.Seconds,
		ex.Speed,
		ex.Grade)
	if err != nil {
		log.Fatal(err)
	}
	return err

}
func (w *Workout) CreateInDB() error {
	timefmt := "2006-01-02T15:04:05"
	db, err := sql.Open("sqlite3", db_path)
	sqlstatement, err := db.Prepare(`
        INSERT INTO workout(date)
        VALUES(?);
        `)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	res, err := sqlstatement.Exec(w.Time.Format(timefmt))
	id, err := res.LastInsertId()
	w.Id = uint64(id)
	fmt.Println("adasdsada")
	fmt.Println(w.Id)
	return err
}
func (w *Workout) SaveSetInDB() error {
	db, err := sql.Open("sqlite3", db_path)
	sqlstatement, err := db.Prepare(`
    INSERT INTO sets(exercise,reps,weight
    ,workout) SELECT ?,?,?, id from workout
    WHERE id=?
    `)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = sqlstatement.Exec(
		w.Sets[len(w.Sets)-1].Exercise,
		w.Sets[len(w.Sets)-1].Reps,
		w.Sets[len(w.Sets)-1].Weight,
		w.Id)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func (w *Workout) AppendSet(s *Set) error {
	w.Sets = append(w.Sets, *s)
	return nil
}
func (w *Workout) FormattedDate() string {
	fmt.Println(w.Time)
	return w.Time.Format("2006-01-02 15:04:05")
}
func GetExercises() []Exercise {
	sqlstatement := "select exercise.id, exercise.name, exercise.description, exercise.reps, exercise.weight, exercise.seconds, exercise.speed, exercise.grade from exercise"
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(sqlstatement)
	if err != nil {
		log.Fatal(err)
	}
	var exercises []Exercise
	for rows.Next() {
		exercise := new(Exercise)
		err := rows.Scan(&exercise.Id,
			&exercise.Name,
			&exercise.Description,
			&exercise.Reps,
			&exercise.Weight,
			&exercise.Seconds,
			&exercise.Speed,
			&exercise.Grade,
		)
		if err != nil {
			log.Fatal(err)
		}
		exercises = append(exercises, *exercise)
	}
	return exercises
}

func LoadWorkout(id uint64) (Workout, error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlstatement := "select workout.date from workout where workout.id =" + fmt.Sprintf("%v", id)

	rows, err := db.Query(sqlstatement)
	if err != nil {
		log.Fatal(err)
	}
	w := Workout{Id: id}
	for rows.Next() {
		err = rows.Scan(&w.Time)
	}
	if err != nil {
		log.Fatal(err)
	}

	sqlstatement = `SELECT sets.id, sets.reps, sets.weight, exercise.name
        FROM sets, exercise
        WHERE exercise.Id = sets.exercise AND sets.workout = ` + fmt.Sprintf("%v", w.Id)
	rows, err = db.Query(sqlstatement)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		s := Set{}
		err := rows.Scan(
			&s.Id,
			&s.Reps,
			&s.Weight,
			&s.Exercise,
		)
		if err != nil {
			log.Fatal(err)
		}
		w.Sets = append(w.Sets, s)
	}
	return w, nil
}

func ExerciseTaskFunc(w http.ResponseWriter, r *http.Request) {
	arg := r.URL.Path[len("/exercise/"):]
	if arg == "create" {
		if r.Method == http.MethodPost {
			s := Exercise{
				Name:        r.FormValue("name"),
				Reps:        r.FormValue("reps") == "on",
				Description: r.FormValue("description"),
				Weight:      r.FormValue("weight") == "on",
				Seconds:     r.FormValue("seconds") == "on",
				Speed:       r.FormValue("speed") == "on",
				Grade:       r.FormValue("grade") == "on",
			}
			fmt.Printf("%+v\n", r)
			s.SaveExercise()
		}
		tmpl := template.Must(template.ParseFiles(
			"templates/base/layout.html",
			"templates/exercisecreate.html",
		))
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		exercises := GetExercises()
		tmpl := template.Must(template.ParseFiles(
			"templates/base/layout.html",
			"templates/exerciselist.html",
		))
		err := tmpl.Execute(w, exercises)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func DashboardTaskFunc(w http.ResponseWriter, r *http.Request) {
	sqlstatement := "select workout.date, workout.id from workout"
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(sqlstatement)
	if err != nil {
		log.Fatal(err)
	}
	var workouts []Workout
	for rows.Next() {
		w := Workout{}
		err := rows.Scan(&w.Time, &w.Id)
		if err != nil {
			log.Fatal(err)
		}
		workouts = append(workouts, w)
	}
	tmpl := template.Must(template.ParseFiles(
		"templates/base/layout.html",
		"templates/dashboard.html",
	))
	err = tmpl.Execute(w, workouts)
	if err != nil {
		log.Fatal(err)
	}
}

func WorkoutTaskFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	idstr := r.URL.Path[len("/workout/"):]
	fmt.Println(idstr)
	var workout Workout
	// If new workout, create it in DB and http redirect to its Id
	if idstr == "" {
		workout.Time = time.Now()
		err := workout.CreateInDB()
		if err != nil {
			log.Fatal(err)
		}
		// Temporary redirect /workout/newid
		http.Redirect(w, r, "/workout/"+fmt.Sprint(workout.Id), 307)
		return
	}

	if idstr != "" {
		id, _ := strconv.ParseUint(idstr, 10, 64)
		workout, _ = LoadWorkout(id)
		fmt.Printf("%+v", workout)
	}
	// Process form input
	if r.Method == http.MethodPost {
		exercise := r.FormValue("exercise")
		reps, _ := strconv.ParseUint(r.FormValue("reps"), 10, 64)
		weight, _ := strconv.ParseFloat(r.FormValue("weight"), 64)
		// Will print ID rather than Exercise name
		s := Set{Exercise: exercise, Reps: reps, Weight: weight}
		fmt.Printf("%+v", workout)
		workout.AppendSet(&s)
		workout.SaveWorkout()
	}

	// Assemble template root struct and execute the template
	var c Context
	c.Exercises = GetExercises()
	c.Workout = &workout
	tmpl := template.Must(template.ParseFiles(
		"templates/base/layout.html",
		"templates/workout.html",
	))
	tmpl.Execute(w, c)

}
