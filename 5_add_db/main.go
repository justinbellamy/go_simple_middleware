package main

import (
	"log"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/pkg/errors"
	"github.com/joho/godotenv"
)

type server struct {
	router *http.ServeMux
	db *Database
}

// init : invoked before main()
func init() {
    // loads values from .env into the system
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

// main : handle any error returned from run()
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// run : retuns any errors to main()
func run() error {
	s := server{router: http.NewServeMux()}
	err := s.setupDatabase()
	if err != nil {
		return errors.Wrap(err, "database setup failed")
	}
	s.routes()
	// run http server
	log.Fatal(http.ListenAndServe(":9000", s.router))
	return nil
}

// initialize database
func (s *server) setupDatabase() error {
	s.db = &Database{
		Driver:   os.Getenv("SQL_DRIVER"),
		Protocol: os.Getenv("SQL_PROTOCOL"),
		Host:     os.Getenv("SQL_HOST"),
		Port:     os.Getenv("SQL_PORT"),
		Name:     os.Getenv("SQL_NAME"),
		User:     os.Getenv("SQL_USER"),
		Password: os.Getenv("SQL_PASS"),
	}
	err := s.db.Open()
	fmt.Println("db opened")
	if err != nil {
		return err 
	}
	return nil
}

// routes
func (s *server) routes() {
	s.router.HandleFunc("/greet", s.handleGreet())
	s.router.HandleFunc("/health", s.HandleHealth())
	s.router.HandleFunc("/dbversion", s.HandleDatabaseVersion())
	s.router.HandleFunc("/", s.adminOnly(s.handleIndex()))
}

// handler for greet
// nesting request and response type inside the handler can make testing easier
func (s *server) handleGreet() http.HandlerFunc {
	fmt.Println("greet handler setting up")
	type request struct {
		Name string
	}
	type response struct {
		Greeting string `json:"greeting"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// get the request and decode the json from the body
		req := request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			panic(err)
		}

		// create the response data "Hello, person!"
		responseData := response{"Hello, " + req.Name + "!"}

		// build the response json
		response, err := json.Marshal(responseData)
		if err != nil {
			panic(err)
		}

		//Set Content-Type header so that clients will know how to read response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//Write json response back to response
		w.Write(response)
	}
}

// handler for health now returns a function that handles the request
// closure allows us to run code before the handler operates (prepare then use)
func (s *server) HandleHealth() http.HandlerFunc {
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

// handler for health now returns a function that handles the request
// closure allows us to run code before the handler operates (prepare then use)
func (s *server) HandleDatabaseVersion() http.HandlerFunc {
	fmt.Println("database version handler setting up")
	body, _ := s.db.Version()
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, body) //use body
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
