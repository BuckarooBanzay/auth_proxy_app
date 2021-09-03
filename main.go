package main

import (
	"auth_proxy/eventemitter"
	"auth_proxy/types"
	"log"
	"net/http"
)

var logger = log.Default()
var loginMessages = make(chan *types.LoginMessage)
var ee = eventemitter.NewEventEmitter()

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/minetest/channel", MinetestEndpoint)
	mux.HandleFunc("/api/login", LoginEndpoint)

	logger.Printf("Listening on port %d", 8080)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
