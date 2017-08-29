package web

import (
	"fmt"
	"net/http"
)

type Server struct {
	*http.ServeMux
}

func NewServer() *Server {
	s := Server{http.NewServeMux()}
	s.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	return &s

}
