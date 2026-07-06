import { showToast } from './toast.js';
import { reRender, reRenderMessages } from '../scripts/_chat.js';
import { handleIncomingTypingEvent } from '../scripts/_chat.js';

class SocketManager {
    // console.log(window.SharedWorker);

    constructor() {
        this.handlers = {};
        this.worker = new SharedWorker('/services/ws-worker.js');
        console.log('new worker')
        
        this.worker.port.start();

        this.worker.port.onmessage = (e) => {
            console.log('[page] got from worker:', e.data); // <-- checkpoint 3
            const workerMsg = e.data;
            console.log(workerMsg)

            if (workerMsg.type === '__logged_out') {
                window.location.href = '/login';
                return;
            }
            if (workerMsg.type !== '__message') {
                return;
            }

            const msg = workerMsg.payload;
            const type = msg.event_type;
            // console.log('[page] dispatching type:', type); // <-- checkpoint 4
            // console.log(this.handlers)
            if (!this.handlers[type]) {
                console.warn("Unhandled WS event:", type, msg);
                return;
            }

            this.handlers[type].forEach(cb => {
                cb(msg.data, msg);
            });
        };

        this.worker.onerror = (e) => {
            console.error('SharedWorker crashed:', e);
        };

        // disconnect worker when closing the tab (does work)
        window.addEventListener('pagehide', () => {
            if (state.isSelfTyping) {
                ws.send({
                    event_type: "typing:stop",
                    data: {
                        userId: window.profile.id,
                        conversationId: state.currentConversationId,
                        receiverId: state.currentReceiverId,
                    }
                });
            }
            // this.worker.port.postMessage({ type: 'disconnect' }); // is handled by send !
        });
    }

    connect() {
        console.log('ws.connect')
        this.worker.port.postMessage({
            type: 'connect',
            wsUri: window.env.wsUri
        });
    }

    send(data) {
        this.worker.port.postMessage({
            type: 'send',
            payload: data
        });
    }

    on(eventType, callback) {
        if (!this.handlers[eventType]) {
            this.handlers[eventType] = [];
        }

        this.handlers[eventType].push(callback);
        // console.log(this.handlers)
    }

    off(eventType, callback) {
        if (!this.handlers[eventType]) return;

        this.handlers[eventType] =
            this.handlers[eventType]
                .filter(cb => cb !== callback);
    }

    disconnect() { // don't know what is used for ?!
        console.log('ws.disconnect')
        // this.worker.port.postMessage({
        //     type: 'disconnect',
        //     wsUri: window.env.wsUri
        // });
        // usually do nothing
        // other tabs may still be using the socket

        // if removing the last tab didnt close the socket !
        // this.socket?.close();
    }

    logout() {
        this.worker.port.postMessage({ type: 'logout' });
    }
}

export const ws = new SocketManager();
export const onlineUsers = new Set();

ws.on("init", (data) => {
    data.forEach(id => onlineUsers.add(id));
});

ws.on("client_connect", (userId) => {
    onlineUsers.add(userId);
    reRender("connect",userId)
});

ws.on("client_disconnect", (userId) => {
    onlineUsers.delete(userId);
    reRender("disconnect",userId)
});

ws.on("new_post",(data)=>{
    //
});

ws.on("new_message", (data) => {
    showToast(data.text, "success");
    reRenderMessages(data)
});


ws.on("message_sent", (data) => {

    console.log("message sent from another tab:", data);
    reRenderMessages({senderId:window.profile.id,...data});
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