package main

import (
	"auth_proxy/types"
	"encoding/json"
	"net/http"
)

type LoginHandler struct {
	username string
	w        http.ResponseWriter
	done     chan bool
}

func (lh *LoginHandler) OnEvent(name string, payload interface{}) {
	if name != "mod-message" {
		// Wrong event
		return
	}

	msg, ok := payload.(*types.ModMessage)
	if !ok {
		// Wrong type
		return
	}

	if msg.Name != lh.username {
		// Wrong username
		return
	}

	lh.w.WriteHeader(200)
	err := json.NewEncoder(lh.w).Encode(msg)
	if err != nil {
		SendError(lh.w, err.Error())
	}

	lh.done <- true
}
