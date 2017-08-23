package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var listen_port string = ":8888"
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

func (w *Workout) FormattedDate() string {
	fmt.Println(w.Time)
	return w.Time.Format("2006-01-02 15:04:05")
}
func (w *Workout) FormatAsMd() string {
	var wlog string
	for i, s := range w.Sets {
		wlog += "* Set: " + fmt.Sprint(i) + "\n"
		wlog += "    * Exercise: " + fmt.Sprint(s.Exercise) + "\n"
		wlog += "    * Reps: " + fmt.Sprint(s.Reps) + "\n"
		wlog += "    * Weight: " + fmt.Sprint(s.Weight) + "\n"
	}
	return wlog
}

func (w *Workout) AppendSet(s *Set) error {
	w.Sets = append(w.Sets, *s)
	return nil
}

func (w *Workout) SaveSetInDB() error {
	db, err := sql.Open("sqlite3", db_path)
	sqlstatement, err := db.Prepare(`
    INSERT INTO sets(exercise,workout,reps,weight
    ,workout) VALUES(?,?,?,(SELECT id from workout
    WHERE date=?))
    `)
	timefmt := "2006-01-02 15:04:05"
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = sqlstatement.Exec(
		w.Sets[len(w.Sets)-1].Exercise,
		w.Sets[len(w.Sets)-1].Reps,
		w.Sets[len(w.Sets)-1].Weight,
		w.Time.Format(timefmt))
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (w *Workout) FormatAsAsciiTable() string {
	exc_max_len := len("Exercise ")
	rep_max_len := len("Reps ")
	wgt_max_len := len("Weight ")
	var Exercises []string
	var Reps []string
	var Weight []string
	for _, s := range w.Sets {
		Exercises = append(Exercises, fmt.Sprint(s.Exercise))
		Reps = append(Reps, fmt.Sprint(s.Reps))
		Weight = append(Weight, fmt.Sprint(s.Weight))
		if len(Exercises[len(Exercises)-1]) > exc_max_len {
			exc_max_len = len(Exercises[len(Exercises)-1])
		}
		if len(Reps[len(Exercises)-1]) > exc_max_len {
			rep_max_len = len(Reps[len(Exercises)-1])
		}
		if len(Weight[len(Exercises)-1]) > exc_max_len {
			wgt_max_len = len(Weight[len(Exercises)-1])
		}
	}
	Table := "| Exercise "
	Table = Table + strings.Repeat(" ", Max(0, exc_max_len-len("Exercise ")))
	Table = Table + " | Reps "
	Table = Table + strings.Repeat(" ", Max(0, rep_max_len-len("Reps ")))
	Table = Table + " | Weight "
	Table = Table + strings.Repeat(" ", Max(0, wgt_max_len-len("Weight ")))
	Table = Table + " |\n"
	Table += strings.Repeat("-", len(Table)-1) + "\n"
	for i, _ := range Exercises {
		str := fmt.Sprint(Exercises[i])
		Table = Table + "| " + str
		Table = Table + strings.Repeat(" ", Max(0, exc_max_len-len(str)))
		str = fmt.Sprint(Reps[i])
		Table = Table + " | " + str
		Table = Table + strings.Repeat(" ", Max(0, rep_max_len-len(str)))
		str = fmt.Sprint(Weight[i])
		Table = Table + " | " + str
		Table = Table + strings.Repeat(" ", Max(0, wgt_max_len-len(str)))
		Table = Table + " |\n"
	}
	return Table
}

func (w *Workout) CreateWorkoutInDB() error {
	timefmt := "2006-01-02T15:04:05"
	db, err := sql.Open("sqlite3", db_path)
	sqlstatement := `
	INSERT INTO workout(date) 
	VALUES(` + w.Time.Format(timefmt) + `);
	SELECT last_insert_rowid()
	`
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(sqlstatement)
	fmt.Println(rows)
	if err != nil {
		log.Fatal(err)
	}
	err = rows.Scan(&w.Id)
	return err
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

func LoadWorkout(id uint64) (Workout, error) {
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println(id)
	sqlstatement := "select workout.date from workout where workout.id =" + fmt.Sprintf("%v", id)
	fmt.Println(sqlstatement)
	rows, err := db.Query(sqlstatement)
	if err != nil {
		log.Fatal(err)
	}
	w := Workout{Id: id}
	err = rows.Scan(&w.Time)
	fmt.Println(w.Time)
	sqlstatement = "select sets.id, sets.exercise, sets.reps, sets.weight, sets.seconds from sets, workout where sets.workout = " + fmt.Sprintf("%v", w.Id)
	fmt.Println(sqlstatement)
	rows, err = db.Query(sqlstatement)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		s := Set{}
		err := rows.Scan(
			&s.Id,
			&s.Exercise,
			&s.Reps,
			&s.Weight,
			&s.Seconds,
		)
		if err != nil {
			log.Fatal(err)
		}
		w.Sets = append(w.Sets, s)
	}
	return w, nil
}

func WorkoutTaskFunc(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Path[len("/workout/"):]
	var workout Workout
	workout.Time = time.Now()
	fmt.Println(idstr)
	fmt.Printf("%+v\n", workout)
	if idstr != "" {
		id, _ := strconv.ParseUint(idstr, 10, 64)
		workout, _ = LoadWorkout(id)
	}
	if r.Method == http.MethodPost {
		fmt.Printf("Here")
		exercise := r.FormValue("exercise")
		reps, _ := strconv.ParseUint(r.FormValue("reps"), 10, 64)
		weight, _ := strconv.ParseFloat(r.FormValue("weight"), 64)
		s := Set{Exercise: exercise, Reps: reps, Weight: weight}
		workout.AppendSet(&s)
		workout.SaveWorkout()
	}
	var c Context
	c.Exercises = GetExercises()
	c.Workout = &workout
	tmpl := template.Must(template.ParseFiles(
		"templates/workout.tmpl",
		"templates/base/header.tmpl",
		"templates/base/footer.tmpl"))
	/*	tmpl := template.Must(template.New(workout.tmpl).Funcs(
		            template.FuncMap{"FormattedDate": Workout.FormattedDate}).ParseFiles(
				"templates/workout.tmpl",
				"templates/base/header.tmpl",
				"templates/base/footer.tmpl"))
	*/
	tmpl.Execute(w, c)
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
		"templates/dashboard.tmpl",
		"templates/base/header.tmpl",
		"templates/base/footer.tmpl"))
	err = tmpl.Execute(w, workouts)
	if err != nil {
		log.Fatal(err)
	}
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
			"templates/exercisecreate.tmpl",
			"templates/base/header.tmpl",
			"templates/base/footer.tmpl"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		exercises := GetExercises()
		tmpl := template.Must(template.ParseFiles(
			"templates/exerciselist.tmpl",
			"templates/base/header.tmpl",
			"templates/base/footer.tmpl"))
		err := tmpl.Execute(w, exercises)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	http.HandleFunc("/workout/", WorkoutTaskFunc)
	http.HandleFunc("/exercise/", ExerciseTaskFunc)
	http.HandleFunc("/dashboard", DashboardTaskFunc)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "dashboard", 301)
	})

	err := http.ListenAndServe(listen_port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
