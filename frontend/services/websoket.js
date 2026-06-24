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


export function showToast(message, type = 'success', duration = 3000) {
    const container = document.getElementById('toast-container');

    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;

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