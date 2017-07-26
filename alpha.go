package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var listen_port string = ":8080"

type Set struct {
	Exercise string
	Reps     uint64
	Weight   float64
	id       uint64
}
type Workout struct {
	Time time.Time
	Sets []Set
}

type Exercise struct {
	Name string
	Id   uint64
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (w *Workout) FormattedDate() string {
	return fmt.Sprintf("%04d-%02d-%02d", w.Time.Year(), w.Time.Month(), w.Time.Day())
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
	db, err := sql.Open("sqlite3", "./alpha.db")
	sqlstatement, err := db.Prepare(`
    INSERT INTO sets(exercise,reps,weight
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
	db, err := sql.Open("sqlite3", "./alpha.db")
	sqlstatement, err := db.Prepare(`
	INSERT INTO workout(date) 
	VALUES(?);
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	res, err := sqlstatement.Exec(
		w.Time.Format(timefmt))
	fmt.Println(res)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (w *Workout) SaveWorkout() error {
	return w.SaveSetInDB()
}

func LoadWorkout(date string) (*Workout, error) {
	sqlstatement := "select sets.id, sets.exercise, sets.reps, sets.weight, sets.seconds from sets, workout where sets.workout = workout.id and workout.date = " + date
	db, err := sql.Open("sqlite3", "./alpha.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Exec(sqlstatement)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows)

	filename := "data/workouts/" + date + ".txt"
	f, err := os.Open(filename)
	defer f.Close()

	layout := "2006-01-02"
	d, _ := time.Parse(layout, date)
	w := Workout{Time: d}
	if err != nil {
		return &w, nil
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		matchedSet, _ := regexp.MatchString("Set: ", scanner.Text())
		matchedExercise, _ := regexp.MatchString("Exercise: ", scanner.Text())
		matchedReps, _ := regexp.MatchString("Reps: ", scanner.Text())
		matchedWeight, _ := regexp.MatchString("Weight: ", scanner.Text())
		if matchedSet {
			w.Sets = append(w.Sets, Set{})
		} else if matchedExercise {
			w.Sets[len(w.Sets)-1].Exercise = strings.Split(scanner.Text(), ":")[1]
		} else if matchedReps {
			w.Sets[len(w.Sets)-1].Reps, _ = strconv.ParseUint(strings.Split(scanner.Text(), ":")[1], 10, 64)
		} else if matchedWeight {
			w.Sets[len(w.Sets)-1].Weight, _ = strconv.ParseFloat(strings.Split(scanner.Text(), ":")[1], 32)
		}
	}
	return &w, nil
}

func WorkoutTaskFunc(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/workout/"):]
	t := time.Now()
	var workout *Workout
	if id != "" {
		layout := "2006-01-02"
		t, _ = time.Parse(layout, id)
	}
	workout, _ = LoadWorkout(fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day()))
	workout.CreateWorkoutInDB()
	if r.Method == http.MethodPost {
		exercise := r.FormValue("exercise")
		reps, _ := strconv.ParseUint(r.FormValue("reps"), 10, 64)
		weight, _ := strconv.ParseFloat(r.FormValue("weight"), 64)
		s := Set{Exercise: exercise, Reps: reps, Weight: weight}
		workout.AppendSet(&s)
		workout.SaveWorkout()

	}
	fmt.Println(workout.FormatAsAsciiTable())
	tmpl := template.Must(template.ParseFiles(
		"templates/workout.tmpl",
		"templates/base/header.tmpl",
		"templates/base/footer.tmpl"))

	tmpl.Execute(w, workout)
}

func DashboardTaskFunc(w http.ResponseWriter, r *http.Request) {
	var date time.Time
	sqlstatement := "select workout.date from workout"
	db, err := sql.Open("sqlite3", "./alpha.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(sqlstatement)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows)
	var workouts []string
	timefmt := "2006-01-02"
	for rows.Next() {
		err := rows.Scan(&date)
		if err != nil {
			log.Fatal(err)
		}
		workouts = append([]string{date.Format(timefmt)}, workouts...)
		fmt.Println(workouts[len(workouts)-1])
	}
	tmpl := template.Must(template.ParseFiles(
		"templates/dashboard.tmpl",
		"templates/base/header.tmpl",
		"templates/base/footer.tmpl"))
	tmpl.Execute(w, workouts)
}

func ExerciseTaskFunc(w http.ResponseWriter, r *http.Request) {
	sqlstatement := "select exercise.id, exercise.name from exercise"
	db, err := sql.Open("sqlite3", "./alpha.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(sqlstatement)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows)
	var exercises []Exercise
	for rows.Next() {
		exercise := new(Exercise)
		err := rows.Scan(&exercise.Id, &exercise.Name)
		if err != nil {
			log.Fatal(err)
		}
		exercises = append(exercises, *exercise)
	}
	tmpl := template.Must(template.ParseFiles(
		"templates/exerciselist.tmpl",
		"templates/base/header.tmpl",
		"templates/base/footer.tmpl"))
	err = tmpl.Execute(w, exercises)
	if err != nil {
		log.Fatal(err)
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
