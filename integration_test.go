package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/egansoft/breezy/routing"
	"github.com/egansoft/breezy/server"
)

const testRoutes = `
/echo/[msg]/3x               $ echo [msg] [msg] [msg]
/count/[flag]/grep/[pattern] $ grep [pattern] | wc -[flag]
/the/static/files            : examples/static
/the/static/files/ignored    $ echo ignored
/the/static/files            $ echo ignored
/error                       $ exit 1
`

func TestIntegration(t *testing.T) {
	s, err := initServer(testRoutes)
	if err != nil {
		t.Fatalf("Couldn't parse routes: %v", err)
	}

	res, code := hitEndpoint(s, http.MethodGet, "/echo/hello/3x", "")
	assertEqual(t, code, http.StatusOK)
	assertEqual(t, res, "hello hello hello")

	res, code = hitEndpoint(s, http.MethodPost, "/count/l/grep/y", "y\nn\ny\nn\nn\nn\ny")
	assertEqual(t, code, http.StatusOK)
	assertEqual(t, res, "3")

	res, code = hitEndpoint(s, http.MethodGet, "/the/static/files/text", "")
	assertEqual(t, code, http.StatusOK)
	assertEqual(t, res, "just some text")

	res, code = hitEndpoint(s, http.MethodGet, "/not/a/real/path", "")
	assertEqual(t, code, http.StatusNotFound)

	res, code = hitEndpoint(s, http.MethodGet, "/echo/notquite", "")
	assertEqual(t, code, http.StatusNotFound)

	res, code = hitEndpoint(s, http.MethodGet, "/the/static/files/ignored", "")
	assertEqual(t, code, http.StatusNotFound)

	res, code = hitEndpoint(s, http.MethodGet, "/error", "")
	assertEqual(t, code, http.StatusInternalServerError)
}

func BenchmarkHandling(b *testing.B) {
	s, _ := initServer(testRoutes)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hitEndpoint(s, http.MethodGet, "/echo/hello/3x", "")
		hitEndpoint(s, http.MethodPost, "/count/l/grep/y", "y\nn\ny\nn\nn\nn\ny")
		hitEndpoint(s, http.MethodGet, "/the/static/files/text", "")
		hitEndpoint(s, http.MethodGet, "/not/a/real/path", "")
		hitEndpoint(s, http.MethodGet, "/error", "")
	}
}

func initServer(routefile string) (*server.Server, error) {
	routes := strings.Split(routefile, "\n")
	router, err := routing.Parse(routes)
	if err != nil {
		return nil, err
	}

	return server.New(router), nil
}

func hitEndpoint(s *server.Server, method, path, data string) (string, int) {
	body := bytes.NewReader([]byte(data))
	request := httptest.NewRequest(method, path, body)
	recorder := httptest.NewRecorder()

	s.ServeHTTP(recorder, request)
	response := recorder.Result()
	code := response.StatusCode
	reader := response.Body

	buf := &bytes.Buffer{}
	io.Copy(buf, reader)
	result := buf.String()
	result = strings.TrimSpace(result)
	return result, code
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	if expected != actual {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Expected %v but got %v on line %v", expected, actual, line)
	}
}
