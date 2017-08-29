package main

import (
	"alpha/web"
	"github.com/stretchr/graceful"
	//	"time"
)

func main() {
	port := "8888"
	s := web.NewServer()
	graceful.Run(":"+port, 0, s)
}
