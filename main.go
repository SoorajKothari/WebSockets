package main

import (
	"github.com/gorilla/websocket"
	"log"
	. "net/http"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func socketHandler(writer ResponseWriter, request *Request) {
	connection, err := upGrader.Upgrade(writer, request, nil)

	if err != nil {
		log.Println("Error upgrading connection to WebSocket:", err)
		return
	}

	defer connection.Close()

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error reading message from WebSocket", err)
			break
		}

		log.Printf("Recieved message  : %s\n", message)

		reply := "Got it " + connection.LocalAddr().String() + " " + string(message)

		err = connection.WriteMessage(messageType, []byte(reply))

		if err != nil {
			log.Println("Error writing message to WebSocket: ", err)
			break
		}
	}

}

func main() {
	HandleFunc("/send", socketHandler)
	ListenAndServe(":8083", nil)
}
