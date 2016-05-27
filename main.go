package main

import (
	"log"
	"net/http"
)

// Hacky for demo purposes for sending changes updates to the server socket
// so the browser gets updated.
// type chan int is inferred from the right-hand side
var changesChannel = make(chan int, 10)

func getChangesChannel() chan int {
	return changesChannel
}

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
