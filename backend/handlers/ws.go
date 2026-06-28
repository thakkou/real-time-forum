package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/utilities"
	"forum/ws"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandlerWs(w http.ResponseWriter, r *http.Request) {
	var userId string

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		http.Error(w, `{"error":"not authenticated"}`, http.StatusUnauthorized)
		return
	}

	id, err := utilities.GetUserIDFromCookie(cookie.Value)
	if err != nil {
		http.Error(w, `{"error":"not authenticated"}`, http.StatusUnauthorized)
		return
	}

	userId = strconv.Itoa(id)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade error:", err)
		return
	}

	client := ws.StoreClient(userId, conn)

	go ws.HandleClient(client)
}

func TestBroadcast(w http.ResponseWriter, r *http.Request) {
	ws.BroadcastExcept("", "test_event", map[string]string{
		"message": "hello everyone 👋",
	})
	w.Write([]byte("sent"))
}
