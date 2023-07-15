package main

import (
	. "WebSockets/src/socket"
	"log"
	. "net/http"
)

func main() {
	go ReplyService()

	HandleFunc("/send", ConnectionHandler)
	err := ListenAndServe(":8083", nil)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
	}
}
