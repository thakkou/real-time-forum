//designe the db schema
//creat the crud for messages
  impliment function esentials
    -SendMessage(senderID, receiverID, text)
    -GetConversation(userID, otherUserID, limit, offset)
    -GetLastMessage(userID)
    
   GET /api/users  (get all users to show in the UI)
    GET /api/users/{id}     (get user profile)
   GET /api/messages/conversations
   post  /api/messages (send new private message)
  body  {
  "receiverId": 5,
  "text": "hello"
}
   GET /api/messages/{user2}?offset=10&limit=10 (get 10 by 10 messages)
   //