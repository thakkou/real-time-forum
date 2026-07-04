// Runs in sepaRATE EXECUTION CONTEXT

const ports = []; // new Set(); -> must change methods and use [...x] to use it ! 
let socket = null;

function broadcast(message) {
    for (let i = ports.length - 1; i >= 0; i--) {
        try {
            ports[i].postMessage(message);
        } catch (e) {
            console.warn('Dead port, removing', e);
            ports.splice(i, 1);
        }
    }
    // ports.forEach(port => {
    //     port.postMessage(message);
    // });
}

function connect(wsUri) {

    if (
        socket &&
        (
            socket.readyState === WebSocket.OPEN ||
            socket.readyState === WebSocket.CONNECTING
        )
    ) {
        console.log('ports')
        return;
    }

    console.log('before creating socket')
    socket = new WebSocket(wsUri);
    console.log('socket')

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
        console.log('[worker] raw frame:', event.data); // <-- checkpoint 1
        let msg;
        try {
            msg = JSON.parse(event.data);
        } catch (e) {
            console.error('Bad WS payload:', event.data);
            return;
        }
        console.log('[worker] parsed, broadcasting to', ports.length, 'ports:', msg); // <-- checkpoint 2
        broadcast({ type: '__message', payload: msg });
    };
}

onconnect = (event) => {
    const port = event.ports[0];
    port.id = crypto.randomUUID();
    ports.push(port);
    port.start();

    port.onmessage = (e) => {
        const msg = e.data;
        switch (msg.type) {
            case 'connect':
                // if (socket && socket.readyState === WebSocket.OPEN) {
                //     console.log(socket)
                //     port.postMessage({ type: '__open' });
                // } else {
                //     console.log('connecting...')
                //     connect(msg.wsUri);
                // }
                // break;
                connect(msg.wsUri);
                break;
            case 'send':
                if (socket && socket.readyState === WebSocket.OPEN) {
                    socket.send(JSON.stringify(msg.payload));
                }
                break;
            case 'disconnect':
                const idx = ports.indexOf(port);
                console.log(ports.splice)
                if (idx !== -1) ports.splice(idx, 1);
                break;
        }
    };
};