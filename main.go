package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"

	"github.com/khunmk/pubsub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var ps = &pubsub.PubSub{}

func autoID() string {
	return uuid.NewV4().String()
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := pubsub.Client{
		Id:         autoID(),
		Connection: conn,
	}

	//add this client into the list
	ps.AddClient(client)

	fmt.Println("New client is connected, total : ", len(ps.Clients))

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Something went wrong ", err)

			ps.RemoveClient(client)

			log.Println("total client and subscription", len(ps.Clients), len(ps.Subscriptions))

			return
		}

		// aMessage := []byte("Hi client, i am server")
		// if err := conn.WriteMessage(messageType, aMessage); err != nil {
		// 	log.Println(err)
		// 	return
		// }

		// fmt.Printf("new message : %s ", p)
		ps.HandleReceiveMessage(client, messageType, p)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// payload := map[string]interface{}{
		// 	"message": "hello go",
		// }

		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(payload)
		http.ServeFile(w, r, "static")
	})

	http.HandleFunc("/ws", websocketHandler)

	fmt.Println("Listen and serve on port :3000")

	http.ListenAndServe(":3000", nil)
}
