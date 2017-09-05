package main

import (
	"alpha/controller"
	"alpha/web"
	"github.com/stretchr/graceful"
)

var port string = "8888"

func main() {
	s := web.NewServer()

	workoutsc := controller.ServiceController{"/workout/", controller.WorkoutTaskFunc}
	workoutsc.Register(s)

	exercisesc := controller.ServiceController{"/exercise/", controller.ExerciseTaskFunc}
	exercisesc.Register(s)

	dashboardsc := controller.ServiceController{"/dashboard/", controller.DashboardTaskFunc}
	dashboardsc.Register(s)

	loginsc := controller.ServiceController{"/login/", controller.LoginGET}
	loginsc.Register(s)

	graceful.Run(":"+port, 0, s)
}
