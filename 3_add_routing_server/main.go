package main

import (
	"fmt"
	"log"
	"net/http"
)

type server struct {
	router *http.ServeMux
}

func main() {
	s := server{router: http.NewServeMux()}
	s.routes()

	// run http server
	log.Fatal(http.ListenAndServe(":9000", s.router))
}

// routes
func (s *server) routes() {
	s.router.HandleFunc("/health", s.handleHealth())
	s.router.HandleFunc("/", s.adminOnly(s.handleIndex()))
}

// handler for health now returns a function that handles the request
// closure allows us to run code before the handler operates (prepare then use)
func (s *server) handleHealth() http.HandlerFunc {
	fmt.Println("health handler setting up")
	body := "OK" //prepare body
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, body) //use body
	}
}

// handler for index now returns a function that handles the request
func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "secret stuff")
	}
}

// middleware example
// visit /?admin=true vs /?admin=false
// takes dependency as an argument, example: (h.HandlerFunc)
func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()

		if qs.Get("admin") != "true" {
			http.NotFound(w, r)
			return
		}

		h(w, r)
	}
}
