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

//1-write a message identifier to triger a function based on message type and extract the data
//2-event name
/*
1-post-notification
2-sent-message-event
3-like comment/post
4-typing in progress
*/
//3-data type i should write
/*
{
}
*/

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

/*this is a structure how to sent the response and for who
-new_posts //all users exept sender {
type:"new post created",
data:{
user created:{
nickname,id
},
postData:{
`post data`}
}}
-commented/liked/newcomment/new message/typing  //all those are sent to one
{
type:"x_to_u",
data:{
user created:{
nickname,id
},
xData:{
`post data`}
}}}

*/
