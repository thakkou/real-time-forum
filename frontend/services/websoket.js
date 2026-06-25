import { Toast } from '../components/Toast.js';

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
            console.error(err);
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
        this.socket?.close();
    }
}


export const ws = new WSService();

// ===================================================

export function showToast(message, type = 'success', duration = 3000) {
    const container = document.getElementById('toast-container');

    const toast = Toast(message, type);

    container.appendChild(toast);

    requestAnimationFrame(() => {
        toast.classList.add('show');
    });

    setTimeout(() => {
        toast.classList.remove('show');

        setTimeout(() => {
            toast.remove();
        }, 300);
    }, duration);
}

// showToast('Message received!');
// showToast('User connected', 'success');
// showToast('Connection lost', 'error');
// showToast('New notification', 'warning', 5000);

// socket.onmessage = (event) => {
//     showToast('New message received');
// };

// ===========================================================


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