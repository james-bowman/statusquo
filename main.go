package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func staticHandler(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path == "" {
		path = "index.html"
	}
	if bs, err := Asset(path); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		log.Printf("Error serving '%s': %s\n", path, err.Error())
	} else {
		var reader = bytes.NewBuffer(bs)
		io.Copy(rw, reader)
	}
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(portfolio.checks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Printf("Error occured processing request '%s' %v\n", req.URL.Path, err)
	}
}

func main() {
	port := 3000

	portfolio = NewPortfolio("checks.json")
	portfolio.Start()

	http.Handle("/", http.StripPrefix("/", http.HandlerFunc(staticHandler)))
	http.Handle("/api/checks/", http.HandlerFunc(apiHandler))
	log.Printf("Listening on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
