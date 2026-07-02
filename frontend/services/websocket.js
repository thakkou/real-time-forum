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

            if (workerMsg.type !== '__message') {
                return;
            }

            const msg = workerMsg.payload;
            const type = msg.event_type;
            console.log('[page] dispatching type:', type); // <-- checkpoint 4
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

        // disconnect worker when closing the tab
        window.addEventListener('pagehide', () => {
            this.worker.port.postMessage({ type: 'disconnect' });
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

    disconnect() {
        // usually do nothing
        // other tabs may still be using the socket

        // if removing the last tab didnt close the socket !
        // this.socket?.close();
    }
}

export const ws = new SocketManager();