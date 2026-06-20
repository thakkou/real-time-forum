1-store all conected clients
2-write event of write messages


--actions to do
-start by storing onlie from oflines 
-a message notification 
-post notification 
-react to ur post
-typing in progresse

//start implimenting ws
what i should impliment ??
//1-add a part i store all conected user
map[userID]*ws.Conn
-ws handlers
  -connect
  -read messages
  -route events

-event Router
 "send_message","typing","read"


 --------------------- workflow -------------
 # connct ws
 const ws = new WebSocket("ws://localhost:8080/ws")

 # identify user
 {
  "event": "auth",
  "token": "session_id"
}
# STEP 3 — Send message event
{
  "event": "send_message",
  "data": {
    "conversation_id": 12,
    "text": "hello"
  }
}

🧱 4. Basic workflow (what you should implement)
STEP 1 — Connect WS

Frontend:

const ws = new WebSocket("ws://localhost:8080/ws")

Backend:

upgrade HTTP → WS
STEP 2 — Identify user

After connect:

{
  "event": "auth",
  "token": "session_id"
}

Backend:

validate session
store connection in map
STEP 3 — Send message event

Frontend:

{
  "event": "send_message",
  "data": {
    "conversation_id": 12,
    "text": "hello"
  }
}

Backend:

save message in DB
find receiver
push event
STEP 4 — Receive real-time message

Backend sends:

{
  "event": "new_message",
  "data": {
    "conversation_id": 12,
    "sender_id": 1,
    "text": "hello",
    "created_at": "..."
  }
}

Frontend:

instantly append to chat UI
🔥 5. Events you SHOULD implement for your project

You don’t need everything. This is enough:

💬 Chat events
send_message
new_message
load_messages
👤 User status
user_online
user_offline
⌨️ UX events
typing
stop_typing
📡 system
auth
ping/pong