package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"forum/utilities"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn   *websocket.Conn
	isAuth bool
	id     string
}

var (
	Clients = make(map[string]*Client)
	mu      sync.RWMutex
)

func HandlerWs(w http.ResponseWriter, r *http.Request) {
	var userId string
	isLoggedIn := false

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		id, err := utilities.GetUserIDFromCookie(cookie.Value)
		if err == nil {
			userId = strconv.Itoa(id)
			isLoggedIn = true
		}
	}

	if !isLoggedIn {
		userId = "guest_" + uuid.NewString()
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade error:", err)
		return
	}

	client := &Client{
		conn:   conn,
		isAuth: isLoggedIn,
		id:     userId,
	}

	mu.Lock()
	Clients[userId] = client
	mu.Unlock()

	fmt.Println("[WS] connected:", userId)
}

// here i will add a step to identify users
// if user login i store here token or userId
// depend on token i will run events

//for all clients
/*
  -store it in Clients  done
  -post notification
*/

//for who have the tokens
/*
   -store it in Connected  done
   -message notification send and rescive
   - x react to ur notification
   -typing in progress
*/
