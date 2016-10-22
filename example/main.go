package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/piotrkowalczuk/rmux"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	mux := rmux.NewServeMux(rmux.ServeMuxOpts{
		Interceptor: func(rw http.ResponseWriter, r *http.Request, h http.Handler) {
			start := time.Now()
			h.ServeHTTP(rw, r)
			logger.Printf("request handled: %s in %v", r.URL.String(), time.Since(start))
		},
	})
	mux.Handle("GET/me", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		io.WriteString(rw, `{"me": "some info"}`)
	}))
	mux.Handle("GET/not-me", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusForbidden)
	}))

	logger.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))
}
