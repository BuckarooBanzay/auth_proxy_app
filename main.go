package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var logger = log.Default()

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

var toMt = make(chan json.RawMessage)
var fromMt = make(chan json.RawMessage)

func SendError(w http.ResponseWriter, msg string) {
	logger.Printf("Error: %s", msg)
	w.WriteHeader(500)
	w.Write([]byte(msg))
}

//TODO: check for race conditions
func MinetestEndpoint(w http.ResponseWriter, req *http.Request) {
	logger.Printf("Got mod-request from %s, method: %s", req.Host, req.Method)

	if req.Method == http.MethodGet {
		select {
		case msg := <-toMt:
			logger.Printf("Relaying message to mod")
			err := json.NewEncoder(w).Encode(msg)
			if err != nil {
				SendError(w, err.Error())
				return
			}
		case <-time.After(30 * time.Second):
			// timed out without data
			w.WriteHeader(204)
		}
	} else if req.Method == http.MethodPost {
		msg := json.RawMessage{}
		err := json.NewDecoder(req.Body).Decode(&msg)
		if err != nil {
			SendError(w, err.Error())
			return
		}
		logger.Printf("Received message from mod: %s", string(msg))
		select {
		case fromMt <- msg:
		default:
			SendError(w, "no receiver available")
		}
	}
}

func LoginEndpoint(w http.ResponseWriter, req *http.Request) {
	logger.Printf("Got request from %s", req.Host)

	msg := json.RawMessage{}
	err := json.NewDecoder(req.Body).Decode(&msg)
	if err != nil {
		SendError(w, err.Error())
		return
	}

	select {
	case toMt <- msg:
	case <-time.After(5 * time.Second):
		SendError(w, "timed out sending")
		return
	}

	select {
	case msg = <-fromMt:
		w.WriteHeader(200)
		err = json.NewEncoder(w).Encode(msg)
		if err != nil {
			SendError(w, err.Error())
		}
	case <-time.After(5 * time.Second):
		SendError(w, "timed out receiving")
	}
}
