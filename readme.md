# real-time-forum

## Description

---

## Requirements

* using JS, private messages, real time actions.

* focus on few points:
Registration and Login
Creation of posts
    Commenting posts
Private Messages

* golang handle Websockets (Backend)

* js handle all the Frontend events and clients Websockets

* You will have only one HTML file, so every change of page you want to do, should be handled in the Javascript. This can be called having a single page application.

* use Gorilla websocket: https://pkg.go.dev/github.com/gorilla/websocket

* where to use go routines and channels ?!

* Registration and Login:
    - non logged-in users will only see the registration or login page.
    - register form required fields (Nickname, Age, Gender, First Name, Last Name, E-mail, Password)
    - must be able to connect using either the nickname or the e-mail combined with the password.
    - must be able to log out from any page on the forum.

* Posts and Comments:
    - must be able to:
        + create posts that will have categories.
        + create comments on the posts.
        + see posts in a feed display (see comments only if they click on a post).

* Private Messages:
    - users will be able to send private messages to each other, so you will need to create a chat.
    - the chat will have:
        + A section to show who is online/offline and able to talk to:
            . must be organized by the last message sent (just like discord).
            . If the user is new and does not present messages you must organize it in alphabetic order.
            . must be able to send private messages to the users who are online ?!
            . This section must be visible at all times.
        + A section that when clicked on the user that you want to send a message:
            . reloads the past messages.
            . chats between users must be visible (be able to see the previous messages that you had with the user)
            . chats must reload the last 10 messages and when scrolled up to see more messages you must provide the user with 10 more, without spamming the scroll event.
            . Do not forget what you learned!! (Throttle, Debounce)
        + Messages must have a specific format:
            . A date that shows when the message was sent.
            . The user name, that identifies the user that sent the message.
    - the messages should work in real time (if a user sends a message, the other user should receive the notification of the new message without refreshing the page).

---

## Authors
