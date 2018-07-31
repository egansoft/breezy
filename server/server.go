package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"

	"github.com/egansoft/breezy/config"
	"github.com/egansoft/breezy/routing"
	"github.com/egansoft/breezy/utils"
)

type Server struct {
	httpServer *http.Server
	router     *routing.Router
}

func New(port uint, router *routing.Router) *Server {
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
	b := &bytes.Buffer{}

	writeDebug(b, "You hit: %s\n\n", r.URL.Path)

	path := utils.UrlToPath(r.URL.Path)
	match := s.router.Match(path)
	writeDebug(b, "Match result:\n%v\n\n", match)

	writeDebug(b, "Router tree:\n%s\n\n", s.router.String())

	if match == nil || match.Action == nil {
		fmt.Fprintf(b, "Page not found")
		writeResponse(w, b, http.StatusNotFound)
		return
	}

	ctype := mime.TypeByExtension(r.URL.Path)
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	writeDebug(b, "Output:\n")
	status, err := match.Action.Handle(b, match.Vars, match.Residual)
	if err != nil {
		log.Println(err)
	}

	writeResponse(w, b, status)
}

func writeResponse(w http.ResponseWriter, b *bytes.Buffer, status int) {
	contentLen, err := b.WriteTo(w)
	if err != nil {
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%v", contentLen))
}

func writeDebug(w io.Writer, msg string, args ...interface{}) {
	if config.DebugMode {
		fmt.Fprintf(w, msg, args...)
	}
}
