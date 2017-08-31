package web

import (
	"net/http"
)

type Server struct {
	*http.ServeMux
}

func NewServer() *Server {
	s := Server{http.NewServeMux()}

	s.HandleFunc("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP)

	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "dashboard", 301)
	})

	return &s

}
