
package main

import (
	"net/http"
	"log"
)

// Hacky for demo purposes for sending changes updates to the server socket
// so the browser gets updated
var changesChannel chan int = make(chan int, 10)
func getChangesChannel() (chan int) {
	return changesChannel
}

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

