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
		qs := req.URL.Query()

		// visit /?admin=true vs /?admin=false
		if qs.Get("admin") != "true" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "secret stuff")
	}

	// routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", indexHandler)

	// run http server
	log.Fatal(http.ListenAndServe(":9000", nil))
}
