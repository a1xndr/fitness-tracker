package main

import (
	"alpha/controllers"
	"alpha/web"
	"github.com/stretchr/graceful"
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

	graceful.Run(":"+port, 0, s)
}
