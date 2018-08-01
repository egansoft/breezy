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
	debugBuf := &bytes.Buffer{}
	defer flushDebug(debugBuf, w)

	bufferDebug(debugBuf, "\nYou hit: %s\n\n", r.URL.Path)

	path := utils.UrlToPath(r.URL.Path)
	match := s.router.Match(path)

	bufferDebug(debugBuf, "Match result:\n%v\n\n", match)
	bufferDebug(debugBuf, "Router tree:\n%s\n\n", s.router.String())

	if match == nil || match.Action == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctype := mime.TypeByExtension(r.URL.Path)
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	status, err := match.Action.Handle(w, match.Vars, match.Residual)
	if err != nil {
		log.Println(err)
	}

	if status != http.StatusOK {
		w.WriteHeader(status)
	}
}

func bufferDebug(w io.Writer, msg string, args ...interface{}) {
	if !config.DebugMode {
		return
	}
	fmt.Fprintf(w, msg, args...)
}

func flushDebug(b *bytes.Buffer, w io.Writer) {
	if !config.DebugMode {
		return
	}

	_, err := b.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}
