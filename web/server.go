package web

import (
	"net/http"
)

type Server struct {
	*http.ServeMux
}

func NewServer() *Server {
	s := Server{http.NewServeMux()}
	return &s

}
