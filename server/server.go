package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/egansoft/silly/tree"
)

type Server struct {
	httpServer *http.Server
	router     *tree.Tree
}

func New(port uint, router *tree.Tree) *Server {
	addr := fmt.Sprintf(":%v", port)
	s := &Server{
		router: router,
	}
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

	fmt.Fprintf(w, "You hit: %s\n\n", r.URL.Path)

	path := getPath(r.URL.Path)
	match := s.router.Match(path)
	fmt.Fprintf(w, "Match result:\n%v\n\n", match)

	fmt.Fprintf(w, "Router tree:\n%s\n", s.router.String())
}

func getPath(url string) []string {
	url = strings.Trim(url, "/")
	return strings.Split(url, "/")
}

func writeError(w http.ResponseWriter, msg string) {

}
