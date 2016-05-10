
package main

import (
	"net/http"
	"log"
)

const (
	ErrorPrefix  = "e:"
	OutputPrefix = "o:"
)

var changesChannel chan int = make(chan int, 10)

func getChangesChannel() (chan int) {
	return changesChannel
}

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

