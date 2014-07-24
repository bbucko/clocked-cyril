package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/bbucko/clocked-cyril/conway"
	"html/template"
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

func createSeed(size int) [][]conway.Cell {
	seed := make([][]conway.Cell, size)
	for i := range seed {
		seed[i] = make([]conway.Cell, size)
	}
	return seed
}

func startGameOfLife(conn *websocket.Conn) {
	seed := createSeed(10)
	seed[4][5] = 1
	seed[5][5] = 1
	seed[6][5] = 1

	gol := new(conway.Board)
	gol.InitWithSeed(10, seed)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				gol.Reaper()
				err := conn.WriteJSON(gol)
				if err != nil {
					log.Println("ticker channel", err)
					return
				}
			}
		}
	}()
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
		t, _ := template.ParseFiles("./web/index2.html")
		t.Execute(w, p)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		go startGameOfLife(conn)
	})

	http.Handle("/js/", http.FileServer(http.Dir("./web/")))
	http.Handle("/img/", http.FileServer(http.Dir("./web/")))

	bind := fmt.Sprintf("%s:%s", host, port)
	log.Println("Starting server on", bind, "with websocket on port", wsPort)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
