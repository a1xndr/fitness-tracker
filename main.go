package main

import (
	"alpha/controllers"
	"alpha/web"
	"github.com/stretchr/graceful"
	"net/http"
)

func main() {
	port := "8888"
	s := web.NewServer()

	workoutsc := controllers.ServiceController{"/workout/", controllers.WorkoutTaskFunc}
	workoutsc.Register(s)

	exercisesc := controllers.ServiceController{"/exercise/", controllers.ExerciseTaskFunc}
	exercisesc.Register(s)

	dashboardsc := controllers.ServiceController{"/dashboard/", controllers.DashboardTaskFunc}
	dashboardsc.Register(s)

	// css and images
	staticsc := controllers.ServiceController{"/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP}
	staticsc.Register(s)

	// redirect GET / to dashboard
	rootsc := controllers.ServiceController{"/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "dashboard", 301)
	}}
	rootsc.Register(s)

	graceful.Run(":"+port, 0, s)
}
