package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Set struct {
	Exercise string
	Reps     uint64
	Weight   float64
}
type Workout struct {
	Time time.Time
	Sets []Set
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (w *Workout) FormattedDate() string {
	return fmt.Sprintf("%04d-%02d-%02d\n", w.Time.Year(), w.Time.Month(), w.Time.Day())
}
func (w *Workout) FormatAsMd() string {
	var wlog string
	for i, s := range w.Sets {
		wlog += "* Set: " + fmt.Sprint(i) + "\n"
		wlog += "    * Exercise: " + fmt.Sprint(s.Exercise) + "\n"
		wlog += "    * Reps:" + fmt.Sprint(s.Reps) + "\n"
		wlog += "    * Weight:" + fmt.Sprint(s.Weight) + "\n"
	}
	return wlog
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

func (w *Workout) SaveWorkout() error {
	w.Sets = append(w.Sets, Set{Exercise: "squat", Reps: 5, Weight: 22})
	w.Sets = append(w.Sets, Set{Exercise: "bench", Reps: 5, Weight: 225})
	w.Sets = append(w.Sets, Set{Exercise: "deadlift", Reps: 10, Weight: 225})
	filename := fmt.Sprintf("%04d-%02d-%02d.txt", w.Time.Year(), w.Time.Month(), w.Time.Day())
	title := fmt.Sprintf("Workout on %04d-%02d-%02d\n", w.Time.Year(), w.Time.Month(), w.Time.Day())
	title += strings.Repeat("=", len(title)-1) + "\n\n"
	wlog := w.FormatAsMd()
	fmt.Println(wlog)
	return ioutil.WriteFile(filename, []byte(title+wlog), 0600)
}

func LoadWorkout(date string) (*Workout, error) {
	//filename := fmt.Sprintf("%04d-%02d-%02d.txt", date.Year(), .Month(), w.Time.Day())
	filename := date + ".txt"
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
		matchedSet, _ := regexp.MatchString("Set:", scanner.Text())
		matchedExercise, _ := regexp.MatchString("Exercise:", scanner.Text())
		matchedReps, _ := regexp.MatchString("Reps:", scanner.Text())
		matchedWeight, _ := regexp.MatchString("Weight:", scanner.Text())
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
	if id != "" {
	}
	workout, _ := LoadWorkout(fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day()))
	tmpl := template.Must(template.ParseFiles("workout.gtpl"))
	tmpl.Execute(w, workout)

}
func ProcessWorkoutFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("workout.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseForm()
		fmt.Println("reps: ", r.Form["reps"])
		fmt.Fprintf(w, "reps: ", r.Form["reps"])
		fmt.Println("weight: ", r.Form["weight"])
		fmt.Fprintf(w, "weight: ", r.Form["weight"])
	}
}

func DeleteTaskFunc(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/delete/"):]
	if id == "all" {
		fmt.Fprintf(w, "Ayo che boi delete all")
	} else {
		id, err := strconv.Atoi(id)
		if err != nil {
			fmt.Fprintf(w, "You fucked up kid.")
		} else {
			fmt.Fprintf(w, "Deleting Task", id)
		}
	}

}

func sayHello(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Whats chice?")
}

func main() {
	http.HandleFunc("/delete/", DeleteTaskFunc)
	http.HandleFunc("/workout/", WorkoutTaskFunc)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
