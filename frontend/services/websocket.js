import { reRender, reRenderMessages, updateTheConv } from '../scripts/_chat.js';
import { showToast } from './toast.js';
import { handleIncomingTypingEvent } from '../scripts/_chat.js';

export const onlineUsers = new Set()

class WSService {
    constructor() {
        this.socket = null;
        this.handlers = {}; // 👈 event registry
    }

    connect() {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            return this.socket;
        }

        this.socket = new WebSocket(window.env.wsUri);

        this.socket.onopen = () => {
            console.log('WS connected');
        };

        this.socket.onclose = () => {
            console.log('WS disconnected');
        };

        this.socket.onerror = (err) => {
            console.log(err);
        };

        this.socket.onmessage = (event) => {
            
            const msg = JSON.parse(event.data);

            const type = msg.event_type;

            if (this.handlers[type]) {
                this.handlers[type].forEach(cb => cb(msg.data, msg));
            } else {
                console.warn("Unhandled WS event:", type, msg);
            }
        };

        return this.socket;
    }

    send(data) {
        if (this.socket?.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(data));
        }
    }

    // 👇 register event listener
    on(eventType, callback) {
        if (!this.handlers[eventType]) {
            this.handlers[eventType] = [];
        }
        this.handlers[eventType].push(callback);
    }

    off(eventType, callback) {
        if (!this.handlers[eventType]) return;

        this.handlers[eventType] = this.handlers[eventType]
            .filter(cb => cb !== callback);
    }

    disconnect() {
        console.log("disconect")
        this.socket?.close();
    }
}


export const ws = new WSService();

ws.on("init", (data) => {
                data.forEach(id => onlineUsers.add(id));
            });


            ws.on("force_logout",(data)=>{
                this.navigate("/login")
            })

            ws.on("client_connect", (userId) => {
                onlineUsers.add(userId);
                reRender("connect",userId)
            });

            ws.on("client_disconnect", (userId) => {
                onlineUsers.delete(userId);
    reRender("disconnect",userId)
            });


            ws.on("new_post",(data)=>{
            });

            ws.on("new_message", (data) => {
                const isMe = data.isMine
                const isNew=data.isNewConversation
               

                    updateTheConv(data,isNew)
                

                 if(!isMe){

              showToast(data.text, "success",data.nickname);

             reRenderMessages(data,false)

                 }else{

                    reRenderMessages(data,true)
                 }


 });

          ws.on("typing:start", (data) => {

  handleIncomingTypingEvent({
    ...data,
    is_typing: true
  });
});

ws.on("typing:stop", (data) => {

  handleIncomingTypingEvent({
    ...data,
    is_typing: false
  });
});


// CORE
// ====

// Web socket manager

// class SocketManager {
//     #ws;

//     connect() {
//         this.#ws = new WebSocket('/ws');

//         this.#ws.onmessage = (event) => {
//             const msg = JSON.parse(event.data);

//             document.dispatchEvent(
//                 new CustomEvent(
//                     `ws:${msg.type}`,
//                     { detail: msg.payload }
//                 )
//             );
//         };
//     }

//     send(type, payload) {
//         this.#ws.send(
//             JSON.stringify({
//                 type,
//                 payload
//             })
//         );
//     }
// }

// export const socket = new SocketManager();