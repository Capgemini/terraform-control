
package main

import (
	"net/http"
	"log"
)

const (
	ErrorPrefix  = "e:"
	OutputPrefix = "o:"
)

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

