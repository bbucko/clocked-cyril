package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Message struct {
	msg string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin : func(r *http.Request) bool {
		return true
	},
}

func echo(conn *websocket.Conn) error {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return nil
		}
		msg := string(p)
		var rsp = ""
		if msg == "Ping" {
			rsp = "Pong"
		} else {
			rsp = msg
		}
		log.Printf("RQ: %s RS: %s", msg, rsp)
		if err = conn.WriteMessage(messageType, []byte(rsp)); err != nil {
			return err
		}
	}
}

func main() {
	log.Println("starting...")

	http.Handle("/", http.FileServer(http.Dir("./web/")))
	http.Handle("/js", http.FileServer(http.Dir("./web/js/")))
	http.Handle("/img", http.FileServer(http.Dir("./web/img/")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
			return
		}

		go echo(conn)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
