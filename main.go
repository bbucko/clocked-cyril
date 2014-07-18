package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	msg string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
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

type Page struct {
	Url    string
	WsPort string
}

func main() {
	log.Println("starting...")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	wsPort := os.Getenv("WSPORT")
	if port == "" {
		port = "8080"
	}

	if wsPort == "" {
		wsPort = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host := strings.Split(r.Host, ":")[0]
		p := &Page{Url: host, WsPort: wsPort}
		t, _ := template.ParseFiles("./web/index.html")
		t.Execute(w, p)
	})

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

	bind := fmt.Sprintf("%s:%s", host, port)
	log.Println("Starting server on", bind, "with websocket on port", wsPort)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
