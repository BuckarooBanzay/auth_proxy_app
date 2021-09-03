package main

import "net/http"

func SendError(w http.ResponseWriter, msg string) {
	logger.Printf("Error: %s", msg)
	w.WriteHeader(500)
	w.Write([]byte(msg))
}
