console.log('Shared worker loaded'); // runs in sepaRATE EXECUTION CONTEXT

const ports = [];

let socket = null;

function broadcast(message) {
    ports.forEach(port => {
        port.postMessage(message);
    });
}

function connect(wsUri) {

    if (
        socket &&
        (
            socket.readyState === WebSocket.OPEN ||
            socket.readyState === WebSocket.CONNECTING
        )
    ) {
        // console.log(ports)
        return;
    }

    socket = new WebSocket(wsUri);

    socket.onopen = () => {
        console.log('Worker WS connected');

        broadcast({
            type: '__open'
        });
    };

    socket.onclose = () => {
        console.log('Worker WS disconnected');

        broadcast({
            type: '__close'
        });

        socket = null;
    };

    socket.onerror = (err) => {
        console.error(err);
    };

    socket.onmessage = (event) => {
console.log('FROM WORKER:', event.data);
        const msg = JSON.parse(event.data);

        broadcast({
            type: '__message',
            payload: msg
        });

    };
}

onconnect = (event) => {

    const port = event.ports[0];

    ports.push(port);
    // console.log(ports)

    port.start();
    port.postMessage('pong');

    port.onmessage = (e) => {

        const msg = e.data;

        switch (msg.type) {

            case 'connect':
                connect(msg.wsUri);
                break;

            case 'send':

                if (
                    socket &&
                    socket.readyState === WebSocket.OPEN
                ) {
                    socket.send(
                        JSON.stringify(msg.payload)
                    );
                }

                break;
        }

    };
};