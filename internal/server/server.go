package server

import (
	"net/http"
	"txparser/internal/parser"

	"github.com/gorilla/mux"
)

type Server struct {
	parser parser.Parser
}

func NewServer(parser parser.Parser) (*Server, error) {
	return &Server{parser: parser}, nil
}

func (s *Server) Start(port string) error {
	r := mux.NewRouter()

	r.HandleFunc("/subscribe", s.HandleSubscribe).Methods("POST")
	r.HandleFunc("/getcurrentblock", s.HandleGetCurrentBlock).Methods("GET")
	r.HandleFunc("/gettransactions", s.HandleGetTransactions).Methods("GET")

	return http.ListenAndServe(port, r)
}
