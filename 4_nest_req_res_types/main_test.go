package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	s := server{router: http.NewServeMux()}
	s.routes()
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.HandleHealth())
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//expected := `{"alive": true}`
	want := `OK`
	if rr.Body.String() != want {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), want)
	}
}

func TestHandleGreet(t *testing.T) {
	s := server{router: http.NewServeMux()}
	s.routes()

	//set up body json
	p := struct {
		Name string `json:"Name"`
	}{
		Name: "Justin",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/greet", &buf)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handleGreet())
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	type response struct {
		Greeting string `json:"greeting"`
	}

	rsp := response{}
	err = json.NewDecoder(rr.Body).Decode(&rsp)
	if err != nil {
		panic(err)
	}

	//expected := `{"alive": true}`
	want := `Hello, Justin!`
	if rsp.Greeting != want {
		t.Errorf("handler returned unexpected body: got %v want %v", rsp.Greeting, want)
	}
}
