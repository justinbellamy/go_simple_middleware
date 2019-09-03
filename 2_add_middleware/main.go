package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	// handler for health
	healthHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("health handler setting up")
		body := "OK"
		fmt.Fprintf(w, body)
	}

	// handler for index
	indexHandler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "secret stuff")
	}

	// routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", adminOnly(indexHandler))

	// run http server
	log.Fatal(http.ListenAndServe(":9000", nil))
}

// middleware example
// visit /?admin=true vs /?admin=false
func adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()

		if qs.Get("admin") != "true" {
			http.NotFound(w, r)
			return
		}

		h(w, r)
	}
}
