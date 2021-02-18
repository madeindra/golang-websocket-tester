package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients []*websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func home(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	clients = append(clients, conn)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", msg)

		err = conn.WriteMessage(msgType, msg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for _, c := range clients {
		c.WriteMessage(1, req)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(req)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
