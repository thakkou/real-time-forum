class SocketManager {
    // console.log(window.SharedWorker);

    constructor() {
        this.handlers = {};
        this.worker = new SharedWorker('/workers/ws-worker.js');

        // try {
        //     console.log("before");
        //     this.worker = new SharedWorker(
        //         '/workers/ws-worker.js'//,
        //         // {
        //         //     type: 'module'
        //         // }
        //     );
        //     console.log(this.worker);
        //     console.log(this.worker.port);
        // } catch (err) {
        //     console.error(err);
        // }
        
        this.worker.port.start();
        // this.worker.port.postMessage('ping');

        this.worker.port.onmessage = (e) => {

            const workerMsg = e.data;

            if (workerMsg.type !== '__message') {
                return;
            }

            const msg = workerMsg.payload;

            const type = msg.event_type;

            if (!this.handlers[type]) {
                return;
            }

            this.handlers[type].forEach(cb => {
                cb(msg.data, msg);
            });

        };

    }

//     connect() {
//         if (this.socket && this.socket.readyState === WebSocket.OPEN) {
//             return this.socket;
//         }

//         this.socket = new WebSocket(window.env.wsUri);

//         this.socket.onopen = () => {
//             console.log('WS connected');
//         };

//         this.socket.onclose = () => {
//             console.log('WS disconnected');
//         };

//         this.socket.onerror = (err) => {
//             console.error(err);
//         };

//         this.socket.onmessage = (event) => {
            
//             const msg = JSON.parse(event.data);

//             const type = msg.event_type;

//             if (this.handlers[type]) {
//                 this.handlers[type].forEach(cb => cb(msg.data, msg));
//             } else {
//                 console.warn("Unhandled WS event:", type, msg);
//             }
//         };

//         return this.socket;
//     }
    connect() {

        this.worker.port.postMessage({
            type: 'connect',
            wsUri: window.env.wsUri
        });

    }

//     send(data) {
//         if (this.socket?.readyState === WebSocket.OPEN) {
//             this.socket.send(JSON.stringify(data));
//         }
//     }
    send(data) {

        this.worker.port.postMessage({
            type: 'send',
            payload: data
        });

    }

//     // 👇 register event listener
//     on(eventType, callback) {
//         if (!this.handlers[eventType]) {
//             this.handlers[eventType] = [];
//         }
//         this.handlers[eventType].push(callback);
//     }
    on(eventType, callback) {
        if (!this.handlers[eventType]) {
            this.handlers[eventType] = [];
        }

        this.handlers[eventType].push(callback);

    }

//     off(eventType, callback) {
//         if (!this.handlers[eventType]) return;

//         this.handlers[eventType] = this.handlers[eventType]
//             .filter(cb => cb !== callback);
//     }
    off(eventType, callback) {

        if (!this.handlers[eventType]) {
            return;
        }

        this.handlers[eventType] =
            this.handlers[eventType]
                .filter(cb => cb !== callback);

    }

    disconnect() {
        // usually do nothing
        // other tabs may still be using the socket

        // if removing the last tab didnt close the socket !
        // this.socket?.close();
    }
}

export const ws = new SocketManager();