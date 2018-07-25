package server

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func New(port uint) *Server {
	addr := fmt.Sprintf(":%v", port)
	s := &Server{}
	httpServer := &http.Server{
		Addr:    addr,
		Handler: s,
	}
	s.httpServer = httpServer
	return s
}

func (s *Server) Start() {
	log.Println("Now serving")
	log.Fatal(s.httpServer.ListenAndServe())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Unsupported method")
		return
	}

	fmt.Fprintf(w, r.URL.Path)
}

func writeError(w http.ResponseWriter, msg string) {

}
