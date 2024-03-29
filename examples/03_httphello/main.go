package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.Handle("POST /add", additionHandler{})
	http.Handle("GET /", http.HandlerFunc(handleRoot))
	http.Handle("GET /other", writingHandler(notFound))
	http.Handle("GET /status", responseStatus(204))

	addr := ":32500"
	err := http.ListenAndServe(addr, nil)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Printf("failed to start a server on addr: %s", addr)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("foo", "bar")
	w.WriteHeader(200)
	_, err := fmt.Fprint(w, "Hello World!")
	if err != nil {
		log.Printf("failed to write response body")
	}
}

type additionHandler struct{}

func (a additionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req additionRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("failed to decode request body: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	sum := 0.0
	for _, addend := range req.Addends {
		sum += addend
	}

	// when status code is not set, it defaults to 200
	err = json.NewEncoder(w).Encode(additionResponseBody{Sum: sum})
	if err != nil {
		log.Printf("failed to write response body")
	}
}

type additionRequestBody struct {
	Addends []float64 `json:"addends"`
}

type additionResponseBody struct {
	Sum float64 `json:"sum"`
}

type writingHandler func(w http.ResponseWriter)

func (h writingHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	h(w)
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

type responseStatus int

func (s responseStatus) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(int(s))
}
