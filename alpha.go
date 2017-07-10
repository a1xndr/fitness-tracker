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
func (w *Workout) SaveWorkout() error {
	w.Sets = append(w.Sets, Set{Exercise: "squat", Reps: 5, Weight: 22})
	w.Sets = append(w.Sets, Set{Exercise: "bench", Reps: 5, Weight: 225})
	w.Sets = append(w.Sets, Set{Exercise: "deadlift", Reps: 10, Weight: 225})
	filename := fmt.Sprintf("%04d-%02d-%02d.txt", w.Time.Year(), w.Time.Month(), w.Time.Day())
	title := fmt.Sprintf("Workout on %04d-%02d-%02d\n", w.Time.Year(), w.Time.Month(), w.Time.Day())
	title += strings.Repeat("=", len(title)-1) + "\n\n"
	var wlog string
	for i, s := range w.Sets {
		wlog += "* Set: " + fmt.Sprint(i) + "\n"
		wlog += "    * Exercise: " + fmt.Sprint(s.Exercise) + "\n"
		wlog += "    * Reps:" + fmt.Sprint(s.Reps) + "\n"
		wlog += "    * Weight:" + fmt.Sprint(s.Weight) + "\n"
	}
	fmt.Println(wlog)
	return ioutil.WriteFile(filename, []byte(title+wlog), 0600)
}

func LoadWorkout(date string) (*Workout, error) {
	//filename := fmt.Sprintf("%04d-%02d-%02d.txt", date.Year(), .Month(), w.Time.Day())
	filename := date + ".txt"
	f, _ := os.Open(filename)
	defer f.Close()

	layout := "2006-01-02"
	d, _ := time.Parse(layout, date)
	w := Workout{Time: d}
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
	if id == "" {
		workout := &Workout{Time: time.Now()}
		workout.SaveWorkout()
	} else {
		workout, _ := LoadWorkout(fmt.Sprintf("%04d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day()))
		fmt.Fprintf(w, "%+v\n", workout)
	}
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
