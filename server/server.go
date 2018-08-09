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

const (
	NotFoundResponse            = "Page not found"
	InternalServerErrorResponse = "Internal server error"
)

func New(router *routing.Router) *Server {
	addr := fmt.Sprintf(":%v", config.Port)
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
	log.Printf("Now serving on port %v", config.Port)
	log.Fatal(s.httpServer.ListenAndServe())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	debugBuf := &bytes.Buffer{}
	defer flushDebug(debugBuf, w)

	bufferDebug(debugBuf, "\nYou hit: %s\n\n", r.URL.Path)

	url := r.URL.Path
	path := utils.UrlToPath(url)
	match := s.router.Match(path)

	bufferDebug(debugBuf, "Match result:\n%v\n\n", match)
	bufferDebug(debugBuf, "Router tree:\n%s\n\n", s.router.Serialize())

	if match == nil || match.Action == nil {
		respondWithError(w, http.StatusNotFound, url)
		return
	}

	ctype := mime.TypeByExtension(r.URL.Path)
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	status, err := match.Action.Handle(w, r.Body, match.Vars, match.Residual)
	if err != nil {
		log.Printf("%v", err)
	}

	if status != http.StatusOK {
		respondWithError(w, status, url)
		return
	}

	logResponse(status, url)
}

func respondWithError(w http.ResponseWriter, status int, url string) {
	w.WriteHeader(status)

	switch status {
	case http.StatusNotFound:
		fmt.Fprintf(w, NotFoundResponse)
	case http.StatusInternalServerError:
		fmt.Fprintf(w, InternalServerErrorResponse)
	}
	logResponse(status, url)
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

func logResponse(status int, url string) {
	log.Printf("%v %s", status, url)
}
