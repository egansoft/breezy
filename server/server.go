package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/egansoft/silly/config"
	"github.com/egansoft/silly/tree"
	"github.com/egansoft/silly/utils"
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

	writeDebug(w, "You hit: %s\n\n", r.URL.Path)

	path := utils.UrlToPath(r.URL.Path)
	match := s.router.Match(path)
	writeDebug(w, "Match result:\n%v\n\n", match)

	writeDebug(w, "Router tree:\n%s\n\n", s.router.String())

	if match != nil && match.Action != nil {
		writeDebug(w, "Output:\n")
		match.Action.Handle(w, match.Vars, match.Residual)
	}
}

func writeDebug(w http.ResponseWriter, msg string, args ...interface{}) {
	if config.DebugMode {
		fmt.Fprintf(w, msg, args...)
	}
}

func writeError(w http.ResponseWriter, msg string) {

}
