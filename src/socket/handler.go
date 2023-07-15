package socket

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	. "net/http"
)

type Message struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Content  string    `json:"content"`
}

type Reply struct {
	Content string `json:"content"`
}

var upGrader = websocket.Upgrader{}
var connectionsStore = make(map[uuid.UUID]*websocket.Conn)

var broadcast = make(chan Message)

func ConnectionHandler(writer ResponseWriter, request *Request) {
	connection, err := upGrader.Upgrade(writer, request, nil)

	if err != nil {
		log.Println("Error upgrading connection to WebSocket:", err)
		return
	}

	defer func(connection *websocket.Conn) {
		err := connection.Close()
		if err != nil {
			log.Println("Error closing connection", err)
		}
	}(connection)

	connectionId, err := uuid.NewRandom()
	connectionsStore[connectionId] = connection

	for {
		var message Message
		err := connection.ReadJSON(&message)
		if err != nil {
			log.Println("Error reading message from WebSocket", err)
			delete(connectionsStore, connectionId)
			break
		}
		message.ID = connectionId

		log.Printf("Recieved message %s from %s using connection %s", message.Content, message.Username, connectionId)
		broadcast <- message
	}

}

func ReplyService() {
	for {
		received := <-broadcast
		conn := connectionsStore[received.ID]
		var response Reply
		response.Content = "Hello from server" //we can also send some user input here
		err := conn.WriteJSON(response)
		if err != nil {
			log.Println("Error while sending response", err)
			break
		}
		log.Printf("Responded Successfully to %s", received.ID)
	}
}
