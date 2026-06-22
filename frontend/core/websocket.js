// Web socket manager

class SocketManager {
    #ws;

    connect() {
        this.#ws = new WebSocket('/ws');

        this.#ws.onmessage = (event) => {
            const msg = JSON.parse(event.data);

            document.dispatchEvent(
                new CustomEvent(
                    `ws:${msg.type}`,
                    { detail: msg.payload }
                )
            );
        };
    }

    send(type, payload) {
        this.#ws.send(
            JSON.stringify({
                type,
                payload
            })
        );
    }
}

export const socket = new SocketManager();