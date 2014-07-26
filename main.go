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
	mgo "gopkg.in/mgo.v2"
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
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				gol.Reaper()

				err := conn.WriteJSON(gol.Cells())
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

func getEnv(key string, defVal string) (env string) {
	env = os.Getenv(key)

	if env == "" {
		env = defVal
	}
	return
}

func main() {
	log.Println("starting...")
	host := getEnv("HOST", "")
	port := getEnv("PORT", "8080")
	wsPort := getEnv("WSPORT", "8080")
	mongoURL := getEnv("MONGOHQ_URL", "localhost")

	session, e := mgo.Dial(mongoURL)
	if e != nil {
		log.Fatal(e)
	}
	log.Println(session.BuildInfo())

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
