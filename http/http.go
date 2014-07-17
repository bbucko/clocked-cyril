package http

import (
	"fmt"
	"net/http"
)

type Display struct{}
type Increment struct{}

var allChannels []chan interface{}

func SetupRoutes() {
	http.HandleFunc("/channel", func(w http.ResponseWriter, r *http.Request) {
		newChannel := make(chan interface{})

		go Handle(newChannel)
		allChannels = append(allChannels, newChannel)
		fmt.Println("No ofw  channels:", len(allChannels))
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < len(allChannels); i++ {
			allChannels[i] <- new(Increment)
		}
	})

	http.HandleFunc("/display", func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < len(allChannels); i++ {
			allChannels[i] <- new(Display)
		}
	})
}

func Handle(channel chan interface{}) {
	idx := 0
	for {
		req := <-channel
		switch t := req.(type) {
		case *Display:
			fmt.Println(idx)
		case *Increment:
			idx = idx + 1
		default:
			fmt.Printf("I don't know what to do with %T", t)
		}
	}
}
