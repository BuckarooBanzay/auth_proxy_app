package main

import (
	"auth_proxy/types"
	"encoding/json"
	"net/http"
	"time"
)

func MinetestEndpoint(w http.ResponseWriter, req *http.Request) {
	logger.Printf("Got mod-request from %s, method: %s", req.Host, req.Method)

	if req.Method == http.MethodGet {
		select {
		case msg := <-loginMessages:
			logger.Printf("Relaying message to mod: name=%s", msg.Username)
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
		msg := &types.ModMessage{}
		err := json.NewDecoder(req.Body).Decode(msg)
		if err != nil {
			SendError(w, err.Error())
			return
		}
		logger.Printf("Received message from mod: name=%s", msg.Name)
		ee.Emit("mod-message", msg)
	}
}

func LoginEndpoint(w http.ResponseWriter, req *http.Request) {
	msg := &types.LoginMessage{}
	err := json.NewDecoder(req.Body).Decode(&msg)
	if err != nil {
		SendError(w, err.Error())
		return
	}

	logger.Printf("Got request from host=%s, username=%s", req.Host, msg.Username)

	done := make(chan bool)

	lh := &LoginHandler{
		done:     done,
		username: msg.Username,
		w:        w,
	}
	ee.AddListener(lh)

	logger.Printf("Sending loginmessage to event channel")
	loginMessages <- msg

	select {
	case <-time.After(5 * time.Second):
		SendError(w, "timed out")
	case <-done:
	}
	ee.RemoveListener(lh)
}
