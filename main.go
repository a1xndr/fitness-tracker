package main

import (
	"alpha/controller"
	"alpha/db"
	"alpha/web"
	"github.com/stretchr/graceful"
)

var port string = "8888"
var db_path string = "./alpha.db"

func main() {
	dbinfo := db.DBinfo{db_path}
	db.Connect(dbinfo)

	s := web.NewServer()

	workoutsc := controller.ServiceController{"/workout/", controller.WorkoutTaskFunc}
	workoutsc.Register(s)

	exercisesc := controller.ServiceController{"/exercise/", controller.ExerciseTaskFunc}
	exercisesc.Register(s)

	dashboardsc := controller.ServiceController{"/dashboard/", controller.DashboardTaskFunc}
	dashboardsc.Register(s)

	loginsc := controller.ServiceController{"/login/", controller.LoginGET}
	loginsc.Register(s)

	registersc := controller.ServiceController{"/register/", controller.RegisterPOST}
	registersc.Register(s)

	graceful.Run(":"+port, 0, s)
}
